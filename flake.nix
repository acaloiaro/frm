{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
    tailwindcss.url = "github:acaloiaro/tailwind-cli-extra";
    ess = {
      url = "github:acaloiaro/ess/v2.13.0";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    templ.url = "github:a-h/templ/v0.3.833";
  };

  outputs = {
    self,
    nixpkgs,
    devenv,
    systems,
    tailwindcss,
    ...
  } @ inputs: let
    pkgs = nixpkgs.legacyPackages."x86_64-linux";
    templ = system: inputs.templ.packages.${system}.templ;
    forEachSystem = nixpkgs.lib.genAttrs (import systems);
  in {
    packages = forEachSystem (system: {
      devenv-up = self.devShells.${system}.default.config.procfileScript;
    });
    devShells = forEachSystem (system: let
      config = self.devShells.${system}.default.config;
      executeEss = ''${inputs.ess.packages.${system}.default}/bin/ess --skip-git-add'';
      postgresPort = 5434;
    in {
      default = devenv.lib.mkShell {
        inherit inputs pkgs;
        modules = [
          {
            languages = {
              nix.enable = true;
              go = {
                enable = true;
                package = pkgs.go_1_23;
              };
            };

            packages = with pkgs; [
              go-migrate
              nixpacks
              reflex
              sqlc
              postgresql
              pre-commit
              (templ system)
              tailwindcss.packages.${system}.default
            ];

            enterShell =
              #bash
              ''
                go install github.com/dmarkham/enumer@latest
                export PGOPTIONS=--search_path=frm
                export DATABASE_URL_NON_PGX="postgres://postgres:postgres@localhost:${toString postgresPort}/frm?sslmode=disable"
                export POSTGRES_URL="$DATABASE_URL_NON_PGX&pool_max_conns=100"
                set +a
                ${executeEss}


                ## Fetch dependenies
                export HTMX_VERSION=2.0.3
                export HYPERSCRIPT_VERSION=0.9.13
                export CHOICES_DOT_JS_VERSION=11.0.2
                export SVG_LOADER_VERSION=1.7.1
                export SORTABLE_VERSION=1.15.6

                if [ ! -f ./static/js/hyperscript.js ]; then
                  curl -sL --verbose "https://unpkg.com/hyperscript.org@$HYPERSCRIPT_VERSION" > ./static/js/hyperscript.js
                fi

                if [ ! -f ./static/js/htmx.js ]; then
                  curl -sL --verbose "https://unpkg.com/htmx.org@$HTMX_VERSION" > ./static/js/htmx.js
                fi

                if [ ! -f ./static/js/htmx-response-targets.js ]; then
                  curl -sL --verbose "https://unpkg.com/htmx.org/dist/ext/response-targets.js" > ./static/js/htmx-response-targets.js
                fi

                if [ ! -f ./static/js/choices.min.js ]; then
                  curl -sL --verbose "https://unpkg.com/choices.js@$CHOICES_DOT_JS_VERSION/public/assets/scripts/choices.min.js" > ./static/js/choices.min.js
                fi

                if [ ! -f ./static/js/svg-loader.min.js ]; then
                  curl -sL --verbose "https://unpkg.com/external-svg-loader@$SVG_LOADER_VERSION/svg-loader.min.js" > ./static/js/svg-loader.min.js
                fi

                if [ ! -f ./static/js/Sortable.min.js ]; then
                  curl -sL --verbose "https://unpkg.com/sortablejs@$SORTABLE_VERSION/Sortable.min.js" > ./static/js/Sortable.min.js
                fi

                if [ ! -f ./static/css/choices.min.css ]; then
                  curl -sL --verbose "https://unpkg.com/choices.js@$CHOICES_DOT_JS_VERSION/public/assets/styles/choices.min.css" > ./static/css/choices.min.css
                fi
                run-show-help
              '';
            process.managers.process-compose.unixSocket.enable = true;

            pre-commit.hooks.env-sample-sync = {
              enable = true;
              always_run = true;
              pass_filenames = false;
              name = "env-sample-sync";
              description = "Sync secrets to env.sample";
              entry = executeEss;
            };

            scripts = with pkgs; {
              run-show-help = {
                description = "Show this help text";
                exec = ''
                  echo
                  echo Helper scripts available:
                  echo
                  ${pkgs.gnused}/bin/sed -e 's| |XX|g' \
                    -e 's|=| |' <<EOF | \
                    ${pkgs.util-linuxMinimal}/bin/column -t | \
                    ${pkgs.gnused}/bin/sed -e 's|XX| |g'
                  ${pkgs.lib.generators.toKeyValue {} (pkgs.lib.mapAttrs (name: value: value.description) config.scripts)}
                  EOF
                  echo
                  echo To start the web server and other jobs, run
                  echo
                  echo "    devenv up"
                  echo
                  echo
                '';
              };
              devdb = {
                exec = "${postgresql}/bin/psql $DATABASE_URL_NON_PGX frm $*";
                description = "Connect to the development database (local)";
              };

              exec-ess = {
                exec = "${inputs.ess.packages.${system}.default}/bin/ess --skip-git-add";
                description = "Execute 'ess' with default parameters";
              };

              frm-dev = {
                description = "Run the development server";
                exec = ''
                  go generate ./... && go run cmd/dev_server/main.go
                '';
              };

              migrate = {
                description = "Use go-migrate to generate migrations";
                exec = "${go-migrate}/bin/migrate -source file://./db/migrations -database $DATABASE_URL_NON_PGX $*";
              };

              run-generate-models = {
                description = "Generate models from SQLc";
                exec = ''
                  ${sqlc}/bin/sqlc generate && echo sqlc generate done
                '';
              };

              jjd = {
                description = "jujutsu diff, specialized for templ projects";
                exec = ''jj diff '~ glob:"**/*_templ.txt" & ~ glob:"**/*_templ.go"' --git'';
              };
            };

            processes.frm-server = {
              exec = ''
                reflex \
                          --start-service \
                          --inverse-regex=testdata \
                          --inverse-regex='_test.go$' \
                          --inverse-regex='^\.devenv' \
                          --inverse-regex='^\.direnv' \
                          --inverse-regex='^vendor' \
                          --inverse-regex='.*_enumer\.go|.*_enumer_.*|.+\.templ|.+frm-dev$' -v \
                          frm-dev
              '';
              process-compose = {
                readiness_probe = {
                  http_get = {
                    host = "127.0.0.1";
                    scheme = "http";
                    path = "/ping";
                    port = "3000";
                    initial_delay_seconds = 5;
                    period_seconds = 2;
                    timeout_seconds = 5;
                    success_threshold = 1;
                    failure_threshold = 3;
                  };
                };

                depends_on = {
                  postgres = {
                    condition = "process_healthy";
                  };
                };
              };
            };

            processes.tailwindcss = {
              exec = ''
                reflex \
                  --start-service \
                  -r '.*tailwind\.css$|.*\.templ$' \
                  --inverse-regex='\.devenv' \
                  --inverse-regex='\.direnv' \
                  -- tailwindcss -i ./ui/css/tailwind.css -o ./static/css/styles.css -c ./ui/tailwind.config.js
              '';
            };

            processes.templ = {
              exec = ''
                templ generate --watch --proxy="http://localhost:3000"
              '';
            };

            processes.sqlc = {
              exec = ''
                reflex \
                  -s \
                  --inverse-regex='^\.devenv' \
                  --inverse-regex='^\.direnv' \
                  -r '.+\.sql$' \
                  -- run-generate-models
              '';
            };

            services = {
              postgres = {
                enable = true;
                package = pkgs.postgresql_16;
                listen_addresses = "127.0.0.1";
                port = postgresPort;
                initialScript = ''
                  CREATE ROLE postgres WITH PASSWORD 'postgres' SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS;
                  CREATE DATABASE frm;
                  CREATE DATABASE frm_test;
                '';
                settings = {
                  max_connections = 250;
                  log_statement = "all";
                };
              };
            };
          }
        ];
      };
    });
  };
}
