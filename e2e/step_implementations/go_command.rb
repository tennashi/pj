class GoCommand
  def GoCommand.finalizer
    proc {
      FileUtils.rm_r @tmp_dir
    }
  end

  def initialize(src_path, cmd_name = nil)
    src_path = Pathname.new(src_path.to_s)

    cmd_name ||= src_path.expand_path.basename.to_s

    @tmp_dir = Pathname.new(Dir.mktmpdir())
    @path = @tmp_dir / cmd_name
    @env = {}

    ObjectSpace.define_finalizer(self, GoCommand.finalizer)

    `go build -o #{@path} #{src_path}`
  end

  def set_env(env)
    @env = env
  end

  def cd(path)
    @wd = path
  end

  def run(*args)
    opts = {}

    Tempfile.create("stdout", @tmp_dir) do |stdout|
      Tempfile.create("stderr", @tmp_dir) do |stderr|
        opts[:chdir] = @wd if @wd
        opts[:out] = stdout
        opts[:err] = stderr

        system(@env, @path.to_s, *args, opts)
        CommandResult.new($?, stdout, stderr)
      end
    end
  end
end

class CommandResult
  def initialize(exit_code, stdout, stderr)
    @exit_code = exit_code

    File.open(stdout.path) do |stdout|
      @raw_out = stdout.read
    end

    File.open(stderr.path) do |stderr|
      @raw_err = stderr.read
    end
  end

  def success?
    @exit_code == 0
  end

  def to_json
    JSON.parse @raw_out
  end
end
