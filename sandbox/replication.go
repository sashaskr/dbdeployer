package sandbox

import (
	"fmt"
	"github.com/datacharmer/dbdeployer/common"
	"os"
)

type Slave struct {
	Node       int
	Port       int
	ServerId   int
	Name       string
	MasterPort int
}

func CreateMasterSlaveReplication(sdef SandboxDef, origin string, nodes int) {

	sdef.ReplOptions = ReplOptions
	vList := VersionToList(sdef.Version)
	rev := vList[2]
	base_port := sdef.Port + MasterSlaveBasePort + (rev * 100)
	if sdef.BasePort > 0 {
		base_port = sdef.BasePort
	}
	base_server_id := 0
	sdef.DirName = "master"
	for check_port := base_port + 1; check_port < base_port+nodes+1; check_port++ {
		CheckPort(sdef.SandboxDir, sdef.InstalledPorts, check_port)
	}

	err := os.Mkdir(sdef.SandboxDir, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sdef.Port = base_port + 1
	sdef.ServerId = (base_server_id + 1) * 100
	sdef.LoadGrants = false
	master_port := sdef.Port
	if nodes < 2 {
		fmt.Println("Can't run replication with less than 2 nodes")
		os.Exit(1)
	}
	slaves := nodes - 1
	var data common.Smap = common.Smap{
		"Copyright":  Copyright,
		"SandboxDir": sdef.SandboxDir,
		"Slaves":     []common.Smap{},
	}

	fmt.Println("Installing and starting master")
	sdef.LoadGrants = true
	sdef.Multi = true
	sdef.Prompt = "master"
	CreateSingleSandbox(sdef, origin)
	for i := 1; i <= slaves; i++ {
		data["Slaves"] = append(data["Slaves"].([]common.Smap), common.Smap{
			"Node":        i,
			"SandboxDir":  sdef.SandboxDir,
			"MasterPort":  master_port,
			"RplUser":     sdef.RplUser,
			"RplPassword": sdef.RplPassword})
		sdef.LoadGrants = false
		sdef.Prompt = fmt.Sprintf("slave%d", i)
		sdef.DirName = fmt.Sprintf("node%d", i)
		sdef.Port = base_port + i + 1
		sdef.ServerId = (base_server_id + i + 1) * 100
		fmt.Printf("Installing and starting slave %d\n", i)
		CreateSingleSandbox(sdef, origin)
		var data_slave common.Smap = common.Smap{
			"Node":       i,
			"SandboxDir": sdef.SandboxDir,
			"Copyright":  Copyright,
		}
		write_script(ReplicationTemplates, fmt.Sprintf("s%d", i), "slave_template", sdef.SandboxDir, data_slave, true)
		write_script(ReplicationTemplates, fmt.Sprintf("n%d", i+1), "slave_template", sdef.SandboxDir, data_slave, true)
	}
	sdef.SBType = "replication-node"
	sb_desc := common.SandboxDescription{
		Basedir: sdef.Basedir + "/" + sdef.Version,
		SBType:  "master-slave",
		Version: sdef.Version,
		Port:    []int{0},
		Nodes:   slaves,
	}
	common.WriteSandboxDescription(sdef.SandboxDir, sb_desc)

	write_script(ReplicationTemplates, "start_all", "start_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "restart_all", "restart_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "status_all", "status_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "test_sb_all", "test_sb_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "stop_all", "stop_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "send_kill_all", "send_kill_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "use_all", "use_all_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "initialize_slaves", "init_slaves_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "check_slaves", "check_slaves_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "m", "master_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "n1", "master_template", sdef.SandboxDir, data, true)
	write_script(ReplicationTemplates, "test_replication", "test_replication_template", sdef.SandboxDir, data, true)
	fmt.Println(sdef.SandboxDir + "/initialize_slaves")
	common.Run_cmd(sdef.SandboxDir + "/initialize_slaves")
	fmt.Printf("Replication directory installed in %s\n", sdef.SandboxDir)
	fmt.Printf("run 'dbdeployer usage multiple' for basic instructions'\n")
}

func CreateReplicationSandbox(sdef SandboxDef, origin string, topology string, nodes int) {

	Basedir := sdef.Basedir + "/" + sdef.Version
	if !common.DirExists(Basedir) {
		fmt.Printf("Base directory %s does not exist\n", Basedir)
		os.Exit(1)
	}

	sandbox_dir := sdef.SandboxDir
	switch topology {
	case "master-slave":
		sdef.SandboxDir += "/" + MasterSlavePrefix + VersionToName(origin)
	case "group":
		if sdef.SinglePrimary {
			sdef.SandboxDir += "/" + GroupSPPrefix + VersionToName(origin)
		} else {
			sdef.SandboxDir += "/" + GroupPrefix + VersionToName(origin)
		}
		if !GreaterOrEqualVersion(sdef.Version, []int{5, 7, 17}) {
			fmt.Println("Group replication requires MySQL 5.7.17 or greater")
			os.Exit(1)
		}
	default:
		fmt.Println("Unrecognized topology. Accepted: 'master-slave', 'group'")
		os.Exit(1)
	}
	if sdef.DirName != "" {
		sdef.SandboxDir = sandbox_dir + "/" + sdef.DirName
	}

	if common.DirExists(sdef.SandboxDir) {
		fmt.Printf("Directory %s already exists\n", sdef.SandboxDir)
		os.Exit(1)
	}

	/*
		err := os.Mkdir(sdef.SandboxDir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/

	switch topology {
	case "master-slave":
		CreateMasterSlaveReplication(sdef, origin, nodes)
	case "group":
		CreateGroupReplication(sdef, origin, nodes)
	}
}
