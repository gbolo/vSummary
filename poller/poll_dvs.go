package poller

import (
	"context"
	"time"

	"github.com/gbolo/vsummary/common"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func (p *Poller) GetDVS() (list []common.VSwitch, err error) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "pollDatastores")

	// Create view for objects
	moType := "DistributedVirtualSwitch"
	m := view.NewManager(p.VmwareClient.Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := m.CreateContainerView(ctx, p.VmwareClient.Client.ServiceContent.RootFolder, []string{moType}, true)
	if err != nil {
		return
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all matching objects
	var molist []mo.DistributedVirtualSwitch
	err = v.Retrieve(
		ctx,
		[]string{moType},
		[]string{"name", "summary", "config"},
		&molist,
	)
	if err != nil {
		return
	}

	// construct the list
	for _, mo := range molist {

		list = append(list, common.VSwitch{
			Name:      mo.Name,
			Moref:     mo.Self.Value,
			Type:      "DVS",
			Ports:     mo.Summary.NumPorts,
			MaxMtu:    int32(common.GetInt(mo, "Config", "MaxMtu")),
			Version:   mo.Summary.ProductInfo.Version,
			VcenterId: v.Client().ServiceContent.About.InstanceUuid,
		})
	}

	log.Infof("poller fetched %d summaries of %s", len(list), moType)
	return

}
