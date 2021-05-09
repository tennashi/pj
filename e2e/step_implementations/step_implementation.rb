require 'test/unit'
include Test::Unit::Assertions

require 'tmpdir'
require 'fileutils'
require 'pathname'
require 'json'
require_relative './go_command'

class PJ
  def initialize()
    @cmd = GoCommand.new('..', 'pj')
    @subcmds = []
  end

  def cd(path)
    @cmd.cd(path)
  end

  def set_env(env)
    @cmd.set_env(env)
  end

  def set_default_opts(opts)
    @default_opts = opts
  end

  def method_missing(method, *args)
    @subcmds << method.to_s

    if args.length == 0
      return self
    end

    last_arg = args.last
    args = build_args(@subcmds, args.slice(0..-2), last_arg) if last_arg.class == Hash
    args = build_args(@subcmds, args) if last_arg.class == String

    result = @cmd.run(*args)
    @subcmds = []

    return result
  end

  def build_args(subcmds, args, opts = {})
    opts = @default_opts.merge opts

    opts['global'].to_a.flatten + subcmds.map{|subcmd| [subcmd, opts[subcmd].to_a]}.flatten + args.to_a
  end
end

@suite_store = Gauge::DataStoreFactory.suite_datastore

before_suite do
  pj = PJ.new
  pj.set_default_opts({'global' => { '-o' => 'json' }})

  @suite_store.put('cmd', pj)
end

@scenario_store = Gauge::DataStoreFactory.scenario_datastore

before_scenario do
  pj = @suite_store.get('cmd')

  tmp_dir = Dir.mktmpdir()
  pj.set_env({'XDG_DATA_HOME' => tmp_dir})
  pj.cd(tmp_dir)

  @scenario_store.put('tmp-dir', tmp_dir)
  @scenario_store.put('projects', {})
end

after_scenario do
  tmp_dir = @scenario_store.get('tmp-dir')
  FileUtils.rm_r tmp_dir
end

step 'Create a project with specifying the project name as <project_name>' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.init project_name

  tmp_dir = @scenario_store.get('tmp-dir')

  projects_store = @scenario_store.get('projects')
  projects_store[project_name] = result.to_json
  @scenario_store.put('projects', projects_store)

  assert_equal(result.success?, true)
  assert_equal(result.to_json, {
    'name' => project_name,
    'workspaces' => [tmp_dir],
    'currentWorkspace' => tmp_dir,
  })
end

step 'Show the details of the project <project_name>' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.get project_name

  tmp_dir = @scenario_store.get('tmp-dir')

  assert_equal(result.success?, true)
  assert_equal(result.to_json, {
    'name' => project_name,
    'workspaces' => [tmp_dir],
    'currentWorkspace' => tmp_dir,
  })
end

step 'Make sure there is only <project_name> in the list of projects' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.list ""

  tmp_dir = @scenario_store.get('tmp-dir')

  assert_equal(result.success?, true)
  assert_equal(result.to_json, [
    {
      'name' => project_name,
      'workspaces' => [tmp_dir],
      'currentWorkspace' => tmp_dir,
    }
  ])
end

step 'Move to the directory that the user wants to specify as the project name' do ||
  pj = @suite_store.get('cmd')
  tmp_dir = @scenario_store.get('tmp-dir')

  pj.cd(tmp_dir)
end

step 'Create a project without specifying a project name' do ||
  pj = @suite_store.get('cmd')
  result = pj.init ""

  tmp_dir = @scenario_store.get('tmp-dir')
  tmp_dir_name = Pathname.new(tmp_dir).basename.to_s

  assert_equal(result.success?, true)
  assert_equal(result.to_json, {
    'name' => tmp_dir_name,
    'workspaces' => [tmp_dir],
    'currentWorkspace' => tmp_dir,
  })
end

step 'Show the details of the project' do ||
  pj = @suite_store.get('cmd')
  tmp_dir = @scenario_store.get('tmp-dir')
  tmp_dir_name = Pathname.new(tmp_dir).basename.to_s

  result = pj.get tmp_dir_name

  assert_equal(result.success?, true)
  assert_equal(result.to_json, {
    'name' => tmp_dir_name,
    'workspaces' => [tmp_dir],
    'currentWorkspace' => tmp_dir,
  })
end

step 'Check that the current project is <project_name>' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.current ""

  tmp_dir = @scenario_store.get('tmp-dir')

  assert_equal(result.success?, true)
  assert_equal(result.to_json, {
    'name' => project_name,
    'workspaces' => [tmp_dir],
    'currentWorkspace' => tmp_dir,
  })
end

step 'Change the current project to <project_name>' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.change project_name

  assert_equal(result.success?, true)
end

step 'Create a new directory with the name <dir_name>' do |dir_name|
  tmp_dir = @scenario_store.get('tmp-dir')
  dir_path = Pathname.new(tmp_dir) / dir_name
  Dir.mkdir dir_path
end

step 'Move to the directory <dir_name>' do |dir_name|
  pj = @suite_store.get('cmd')
  tmp_dir = @scenario_store.get('tmp-dir')
  dir_path = Pathname.new(tmp_dir) / dir_name

  pj.cd(dir_path)
end

step 'Add the current directory as a workspace to the current project' do ||
  pj = @suite_store.get('cmd')
  #result = pj.workspace 'add'
  result = pj.workspace.add ''

  assert_equal(result.success?, true)
end

step 'Check that the current project has <dir_name> as a workspace' do |dir_name|
  pj = @suite_store.get('cmd')
  result = pj.workspace.list ""

  tmp_dir = @scenario_store.get('tmp-dir')
  dir_path = Pathname.new(tmp_dir) / dir_name

  assert_equal(result.success?, true)
  assert_equal(result.to_json.map{|e| e['path']}.include?(dir_path.to_s), true)
end

step 'Check that the current project has only one workspace' do ||
  pj = @suite_store.get('cmd')
  result = pj.current ""

  assert_equal(result.success?, true)
  assert_equal(result.to_json['workspaces'].length, 1)
end

step 'Merge the project <project_name> into the current project' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.merge project_name

  assert_equal(result.success?, true)
end

step 'Check that the current project has all workspaces that <project_name> had' do |project_name|
  pj = @suite_store.get('cmd')
  result = pj.current ""

  projects_store = @scenario_store.get('projects')
  want = projects_store[project_name]['workspaces']

  assert_equal(result.success?, true)
  assert_equal((want - result.to_json['workspaces']).empty?, true)
end
