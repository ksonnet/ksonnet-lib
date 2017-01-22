local core = (import "../../../kube/core.libsonnet");
local kubeUtil = (import "../../../kube/util.libsonnet");

local env = core.v1.env + kubeUtil.app.v1.env;
local metadata = core.v1.metadata;
local port = core.v1.port;
local volume = core.v1.volume;

{
  postgres:: {
    config:: {
      "01_create_mattermost_production.sql" :
        "CREATE DATABASE mattermost_production WITH OWNER gitlab;"
    },

    deploy:: {
      Env(appConfigMapName, secretsConfigMapName)::
        env.array.FromConfigMapName(appConfigMapName, {
          "POSTGRES_USER": "postgres_user",
          "POSTGRES_DB": "postgres_db",
        }) +
        env.array.FromConfigMapName(secretsConfigMapName, {
          "POSTGRES_PASSWORD": "postgres_password",
        }) + [
          env.Variable("DB_EXTENSION", "pg_trgm"),
        ],
    },
  },

  gitlab:: {
    configData:: {
      local gkeDomain = std.extVar("GITLAB_GKE_DOMAIN"),

      "external_scheme": "https",
      "external_hostname": "gitlab.%s" % gkeDomain,
      "registry_external_scheme": "https",
      "registry_external_hostname": "registry.%s" % gkeDomain,
      "mattermost_external_scheme": "https",
      "mattermost_external_hostname": "mattermost.%s" % gkeDomain,
      "mattermost_app_uid": "aadas",
      "postgres_user": "gitlab",
      "postgres_db": "gitlab_production",
    },

    secretsData:: {
      "postgres_password": "NDl1ZjNtenMxcWR6NXZnbw==",
      "initial_shared_runners_registration_token": "NDl1ZjNtenMxcWR6NXZnbw==",
      "mattermost_app_secret": "NDl1ZjNtenMxcWR6NXZnbw==",
    },

    patches:: { "fix-git-hooks.patch": $.gitlab.build.patch },

    deploy:: {
      command: ["/bin/bash", "-c",
        "patch -p1 -d /opt/gitlab/embedded/service/gitlab-rails < /patches/fix-git-hooks.patch && sed -i \"s/environment ({'GITLAB_ROOT_PASSWORD' => initial_root_password }) if initial_root_password/environment ({'GITLAB_ROOT_PASSWORD' => initial_root_password, 'GITLAB_SHARED_RUNNERS_REGISTRATION_TOKEN' => node['gitlab']['gitlab-rails']['initial_shared_runners_registration_token'] })/g\" /opt/gitlab/embedded/cookbooks/gitlab/recipes/database_migrations.rb && exec /assets/wrapper"],
      Env(appConfigMapName, secretsConfigMapName)::
        env.array.FromConfigMapName(appConfigMapName, {
          "GITLAB_EXTERNAL_SCHEME": "external_scheme",
          "GITLAB_EXTERNAL_HOSTNAME": "external_hostname",
          "GITLAB_REGISTRY_EXTERNAL_SCHEME": "registry_external_scheme",
          "GITLAB_REGISTRY_EXTERNAL_HOSTNAME": "registry_external_hostname",
          "GITLAB_MATTERMOST_EXTERNAL_SCHEME": "mattermost_external_scheme",
          "GITLAB_MATTERMOST_EXTERNAL_HOSTNAME": "mattermost_external_hostname",
          "POSTGRES_USER": "postgres_user",
          "POSTGRES_DB": "postgres_db",
          "MATTERMOST_APP_UID": "mattermost_app_uid",
        }) +
        env.array.FromConfigMapName(secretsConfigMapName, {
          "POSTGRES_PASSWORD": "postgres_password",
          "GITLAB_INITIAL_SHARED_RUNNERS_REGISTRATION_TOKEN":
            "initial_shared_runners_registration_token",
          "MATTERMOST_APP_SECRET": "mattermost_app_secret",
        }) + [
          env.Variable(
            "GITLAB_POST_RECONFIGURE_SCRIPT",
            |||
                /opt/gitlab/bin/gitlab-rails runner -e production 'Doorkeeper::Application.where(uid: ENV["MATTERMOST_APP_UID"], secret: ENV["MATTERMOST_APP_SECRET"], redirect_uri: "#{ENV["GITLAB_MATTERMOST_EXTERNAL_SCHEME"]}://#{ENV["GITLAB_MATTERMOST_EXTERNAL_HOSTNAME"]}/signup/gitlab/complete\r\n#{ENV["GITLAB_MATTERMOST_EXTERNAL_SCHEME"]}://#{ENV["GITLAB_MATTERMOST_EXTERNAL_HOSTNAME"]}/login/gitlab/complete", name: "GitLab Mattermost").first_or_create;'
            |||),
          env.Variable(
            "GITLAB_OMNIBUS_CONFIG",
            |||
                external_url "#{ENV['GITLAB_EXTERNAL_SCHEME']}://#{ENV['GITLAB_EXTERNAL_HOSTNAME']}"
                registry_external_url "#{ENV['GITLAB_REGISTRY_EXTERNAL_SCHEME']}://#{ENV['GITLAB_REGISTRY_EXTERNAL_HOSTNAME']}"
                mattermost_external_url "#{ENV['GITLAB_MATTERMOST_EXTERNAL_SCHEME']}://#{ENV['GITLAB_MATTERMOST_EXTERNAL_HOSTNAME']}"

                gitlab_rails['initial_shared_runners_registration_token'] = ENV['GITLAB_INITIAL_SHARED_RUNNERS_REGISTRATION_TOKEN']

                nginx['enable'] = false
                registry_nginx['enable'] = false
                mattermost_nginx['enable'] = false

                gitlab_workhorse['listen_network'] = 'tcp'
                gitlab_workhorse['listen_addr'] = '0.0.0.0:8005'

                mattermost['service_address'] = '0.0.0.0'
                mattermost['service_port'] = '8065'

                registry['registry_http_addr'] = '0.0.0.0:8105'

                postgresql['enable'] = false
                gitlab_rails['db_host'] = 'gitlab-postgresql'
                gitlab_rails['db_password'] = ENV['POSTGRES_PASSWORD']
                gitlab_rails['db_username'] = ENV['POSTGRES_USER']
                gitlab_rails['db_database'] = ENV['POSTGRES_DB']

                redis['enable'] = false
                gitlab_rails['redis_host'] = 'gitlab-redis'

                mattermost['file_directory'] = '/gitlab-data/mattermost';
                mattermost['sql_driver_name'] = 'postgres';
                mattermost['sql_data_source'] = "user=#{ENV['POSTGRES_USER']} host=gitlab-postgresql port=5432 dbname=mattermost_production password=#{ENV['POSTGRES_PASSWORD']} sslmode=disable";
                mattermost['gitlab_enable'] = true;
                mattermost['gitlab_secret'] = ENV['MATTERMOST_APP_SECRET'];
                mattermost['gitlab_id'] = ENV['MATTERMOST_APP_UID'];
                mattermost['gitlab_scope'] = '';
                mattermost['gitlab_auth_endpoint'] = "#{ENV['GITLAB_EXTERNAL_SCHEME']}://#{ENV['GITLAB_EXTERNAL_HOSTNAME']}/oauth/authorize";
                mattermost['gitlab_token_endpoint'] = "#{ENV['GITLAB_EXTERNAL_SCHEME']}://#{ENV['GITLAB_EXTERNAL_HOSTNAME']}/oauth/token";
                mattermost['gitlab_user_api_endpoint'] = "#{ENV['GITLAB_EXTERNAL_SCHEME']}://#{ENV['GITLAB_EXTERNAL_HOSTNAME']}/api/v3/user"

                manage_accounts['enable'] = true
                manage_storage_directories['manage_etc'] = false

                gitlab_shell['auth_file'] = '/gitlab-data/ssh/authorized_keys'
                git_data_dir '/gitlab-data/git-data'
                gitlab_rails['shared_path'] = '/gitlab-data/shared'
                gitlab_rails['uploads_directory'] = '/gitlab-data/uploads'
                gitlab_ci['builds_directory'] = '/gitlab-data/builds'
                gitlab_rails['registry_path'] = '/gitlab-registry'
                gitlab_rails['trusted_proxies'] = ["10.0.0.0/8","172.16.0.0/12","192.168.0.0/16"]

                prometheus['enable'] = true
                node_exporter['enable'] = true
            |||),
        ],
    },

    build:: {
      patch:: |||
      diff --git a/app/models/repository.rb b/app/models/repository.rb
      index 30be726243..0776c7ccc5 100644
      --- a/app/models/repository.rb
      +++ b/app/models/repository.rb
      @@ -160,14 +160,18 @@ class Repository
          tags.find { |tag| tag.name == name }
        end

      -  def add_branch(user, branch_name, target)
      +  def add_branch(user, branch_name, target, with_hooks: true)
          oldrev = Gitlab::Git::BLANK_SHA
          ref    = Gitlab::Git::BRANCH_REF_PREFIX + branch_name
          target = commit(target).try(:id)

          return false unless target

      -    GitHooksService.new.execute(user, path_to_repo, oldrev, target, ref) do
      +    if with_hooks
      +      GitHooksService.new.execute(user, path_to_repo, oldrev, target, ref) do
      +        update_ref!(ref, target, oldrev)
      +      end
      +    else
            update_ref!(ref, target, oldrev)
          end

      diff --git a/app/services/commits/change_service.rb b/app/services/commits/change_service.rb
      index 1c82599c57..2d4c9788d0 100644
      --- a/app/services/commits/change_service.rb
      +++ b/app/services/commits/change_service.rb
      @@ -55,7 +55,7 @@ module Commits
            return success if repository.find_branch(new_branch)

            result = CreateBranchService.new(@project, current_user)
      -                                  .execute(new_branch, @target_branch, source_project: @source_project)
      +                                  .execute(new_branch, @target_branch, source_project: @source_project, with_hooks: false)

            if result[:status] == :error
              raise ChangeError, "There was an error creating the source branch: #{result[:message]}"
      diff --git a/app/services/create_branch_service.rb b/app/services/create_branch_service.rb
      index 757fc35a78..a6a3461e17 100644
      --- a/app/services/create_branch_service.rb
      +++ b/app/services/create_branch_service.rb
      @@ -1,5 +1,5 @@
      class CreateBranchService < BaseService
      -  def execute(branch_name, ref, source_project: @project)
      +  def execute(branch_name, ref, source_project: @project, with_hooks: true)
          valid_branch = Gitlab::GitRefValidator.validate(branch_name)

          unless valid_branch
      @@ -26,7 +26,7 @@ class CreateBranchService < BaseService

                          repository.find_branch(branch_name)
                        else
      -                   repository.add_branch(current_user, branch_name, ref)
      +                   repository.add_branch(current_user, branch_name, ref, with_hooks: with_hooks)
                        end

          if new_branch
      diff --git a/app/services/files/base_service.rb b/app/services/files/base_service.rb
      index 9bd4bd464f..1802b932e0 100644
      --- a/app/services/files/base_service.rb
      +++ b/app/services/files/base_service.rb
      @@ -74,7 +74,7 @@ module Files
          end

          def create_target_branch
      -      result = CreateBranchService.new(project, current_user).execute(@target_branch, @source_branch, source_project: @source_project)
      +      result = CreateBranchService.new(project, current_user).execute(@target_branch, @source_branch, source_project: @source_project, with_hooks: false)

            unless result[:status] == :success
              raise_error("Something went wrong when we tried to create #{@target_branch} for you: #{result[:message]}")
  |||,
      // TODO: The above line `|||` is required to not be indented more than the
      // parent? Check if this is a bug.
    },
  },
}