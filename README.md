# dbdeployer

[DBdeployer](https://github.com/datacharmer/dbdeployer) is a tool that deploys MySQL database servers easily.

[![asciicast](https://asciinema.org/a/160541.png)](https://asciinema.org/a/160541)

## Main operations

With dbdeployer, you can deploy a single sandbox, or many sandboxes  at once, with or without replication.

The main commands are **single**, **replication**, and **multiple**, which work with MySQL tarball that have been unpacked into the _sandbox-binary_ directory (by default, $HOME/opt/mysql.)

To use a tarball, you must first run the **unpack** command, which will unpack the tarball into the right directory.

For example:

    $ dbdeployer --unpack-version=8.0.4 unpack mysql-8.0.4-rc-linux-glibc2.12-x86_64.tar.gz
    Unpacking tarball mysql-8.0.4-rc-linux-glibc2.12-x86_64.tar.gz to $HOME/opt/mysql/8.0.4
    .........100.........200.........292

    $ dbdeployer single 8.0.4
    Database installed in $HOME/sandboxes/msb_8_0_4
    . sandbox server started


The program doesn't have any dependencies. Everything is included in the binary. Calling *dbdeployer* without arguments or with '--help' will show the main help screen.

    $ dbdeployer -h
    Makes MySQL server installation an easy task.
    Runs single, multiple, and replicated sandboxes.
    
    Usage:
      dbdeployer [command]
    
    Available Commands:
      delete      delete an installed sandbox
      help        Help about any command
      multiple    create multiple sandbox
      replication create replication sandbox
      sandboxes   List installed sandboxes
      single      deploys a single sandbox
      templates   Admin operations on templates
      unpack      unpack a tarball into the binary directory
      usage       Shows usage of installed sandboxes
      versions    List available versions
    
    Flags:
          --base-port int              Overrides default base-port (for multiple sandboxes)
          --bind-address string        defines the database bind-address  (default "127.0.0.1")
          --config string              config file (default "./dbdeployer.json")
      -p, --db-password string         database password (default "msandbox")
      -u, --db-user string             database user (default "msandbox")
          --gtid                       enables GTID
      -h, --help                       help for dbdeployer
      -i, --init-options strings       mysqld options to run during initialization
          --keep-auth-plugin           in 8.0.4+, does not change the auth plugin
      -c, --my-cnf-options strings     mysqld options to add to my.sandbox.cnf
          --port int                   Overrides default port
          --remote-access string       defines the database access  (default "127.%")
          --rpl-password string        replication password (default "rsandbox")
          --rpl-user string            replication user (default "rsandbox")
          --sandbox-binary string      Binary repository (default "$HOME/opt/mysql")
          --sandbox-directory string   Changes the default sandbox directory
          --sandbox-home string        Sandbox deployment direcory (default "$HOME/sandboxes")
          --use-template strings       [template_name:file_name] Replace existing template with one from file
          --version                    version for dbdeployer
    
    Use "dbdeployer [command] --help" for more information about a command.

The flags listed in the main screen can be used with any commands.
The flags _--my-cnf-options_ and _--init-options_ can be used several times.

If you don't have any tarballs installed in your system, you should first *unpack* it (see an example above).

    $ dbdeployer unpack -h
    If you want to create a sandbox from a tarball, you first need to unpack it
    into the sandbox-binary directory. This command carries out that task, so that afterwards 
    you can call 'single', 'multiple', and 'replication' commands with only the MySQL version
    for that tarball.
    
    Usage:
      dbdeployer unpack MySQL-tarball [flags]
    
    Examples:
    
        $ dbdeployer --unpack-version=8.0.4 unpack mysql-8.0.4-rc-linux-glibc2.12-x86_64.tar.gz
        Unpacking tarball mysql-8.0.4-rc-linux-glibc2.12-x86_64.tar.gz to $HOME/opt/mysql/8.0.4
        .........100.........200.........292
    	
    
    Flags:
      -h, --help                    help for unpack
          --prefix string           Prefix for the final expanded directory
          --unpack-version string   which version is contained in the tarball (mandatory)
    

The main command is *single*, which installs a single sandbox.

    $ dbdeployer single -h
    single installs a sandbox and creates useful scripts for its use.
    MySQL-Version is in the format x.x.xx, and it refers to a directory named after the version
    containing an unpacked tarball. The place where these directories are found is defined by 
    --sandbox-binary (default: $HOME/opt/mysql.)
    For example:
    	dbdeployer single 5.7.21
    
    For this command to work, there must be a directory $HOME/opt/mysql/5.7.21, containing
    the binary files from mysql-5.7.21-$YOUR_OS-x86_64.tar.gz
    Use the "unpack" command to get the tarball into the right directory.
    
    Usage:
      dbdeployer single MySQL-Version [flags]
    
    Flags:
      -h, --help     help for single
          --master   Make the server replication ready
    

If you want more than one sandbox of the same version, without any replication relationship, use the *multiple* command with an optional "--node" flag (default: 3).

    $ dbdeployer multiple -h
    Creates several sandboxes of the same version,
    without any replication relationship.
    For this command to work, there must be a directory $HOME/opt/mysql/5.7.21, containing
    the binary files from mysql-5.7.21-$YOUR_OS-x86_64.tar.gz
    Use the "unpack" command to get the tarball into the right directory.
    
    Usage:
      dbdeployer multiple MySQL-Version [flags]
    
    Examples:
    
    	$ dbdeployer multiple 5.7.21
    	
    
    Flags:
      -h, --help        help for multiple
      -n, --nodes int   How many nodes will be installed (default 3)
    

The *replication* command will install a master and two or more slaves, with replication started. You can change the topology to "group" and get three nodes in peer replication.

    $ dbdeployer replication -h
    The replication command allows you to deploy several nodes in replication.
    Allowed topologies are "master-slave" and "group" (requires 5.7.17+)
    For this command to work, there must be a directory $HOME/opt/mysql/5.7.21, containing
    the binary files from mysql-5.7.21-$YOUR_OS-x86_64.tar.gz
    Use the "unpack" command to get the tarball into the right directory.
    
    Usage:
      dbdeployer replication MySQL-Version [flags]
    
    Examples:
    
    		$ dbdeployer replication 5.7.21
    		# (implies topology = master-slave)
    
    		$ dbdeployer --topology=master-slave replication 5.7.21
    		# (explicitly setting topology)
    
    		$ dbdeployer --topology=group replication 5.7.21
    		$ dbdeployer --topology=group replication 8.0.4 --single-primary
    	
    
    Flags:
      -h, --help              help for replication
      -n, --nodes int         How many nodes will be installed (default 3)
          --single-primary    Using single primary for group replication
      -t, --topology string   Which topology will be installed (default "master-slave")
    

## Multiple sandboxes, same version and type

If you want to deploy several instances of the same version and the same type (for example two single sandboxes of 8.0.4, or two group replication instances with different single-primary setting) you can specify the data directory name and the ports manually.

    $ dbdeployer single 8.0.4
    # will deploy in msb_8_0_4 using port 8004

    $ dbdeployer single 8.0.4 --sandbox-directory=msb2_8_0_4 --port=8005
    # will deploy in msb2_8_0_4 using port 8005

    $ dbdeployer replication 8.0.4 --sandbox-directory=rsandbox22_8_0_4 --base-port=18600
    # will deploy replication in rsandbox22_8_0_4 using ports 18601, 18602, 18603

## Sandbox management

You can list the available MySQL versions with

    $ dbdeployer versions

Also "available" is a recognized alias for this command.

And you can list which sandboxes were already installed

    $ dbdeployer installed  # Aliases: sandboxes, deployed

The command "usage" shows how to use the scripts that were installed with each sandbox.

    $ dbdeployer usage
    
    	USING A SANDBOX
    
    Change directory to the newly created one (default: $SANDBOX_HOME/msb_VERSION 
    for single sandboxes)
    [ $SANDBOX_HOME = $HOME/sandboxes unless modified with flag --sandbox-home ]
    
    The sandbox directory of the instance you just created contains some handy 
    scripts to manage your server easily and in isolation.
    
    "./start", "./status", "./restart", and "./stop" do what their name suggests. 
    start and restart accept parameters that are eventually passed to the server. 
    e.g.:
    
      ./start --server-id=1001
    
      ./restart --event-scheduler=disabled
    
    "./use" calls the command line client with the appropriate parameters,
    Example:
    
        ./use -BN -e "select @@server_id"
    	./use -u root
    
    "./clear" stops the server and removes everything from the data directory, 
    letting you ready to start from scratch. (Warning! It's irreversible!)
    
     USING MULTIPLE SERVER SANDBOX
    On a replication sandbox, you have the same commands (run "dbdeployer usage single"), 
    with an "_all" suffix, meaning that you propagate the command to all the members. 
    Then you have "./m" as a shortcut to use the master, "./s1" and "./s2" to access 
    the slaves (and "s3", "s4" ... if you define more).
    
    In group sandboxes without a master slave relationship (group replication and 
    multiple sandboxes) the nodes can be accessed by ./n1, ./n2, ./n3, and so on.
    
    start_all
    status_all
    restart_all
    stop_all
    use_all
    clear_all
    m
    s1, s2, n1, n2
    
    The scripts "check_slaves" or "check_nodes" give the status of replication in the sandbox.
    

