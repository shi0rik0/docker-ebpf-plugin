// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || arm || arm64 || loong64 || mips64le || mipsle || ppc64le || riscv64

package ebpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadTc returns the embedded CollectionSpec for tc.
func loadTc() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_TcBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load tc: %w", err)
	}

	return spec, err
}

// loadTcObjects loads tc and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*tcObjects
//	*tcPrograms
//	*tcMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadTcObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadTc()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// tcSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type tcSpecs struct {
	tcProgramSpecs
	tcMapSpecs
}

// tcSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type tcProgramSpecs struct {
	TcIngress *ebpf.ProgramSpec `ebpf:"tc_ingress"`
}

// tcMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type tcMapSpecs struct {
	IpIfindexMap *ebpf.MapSpec `ebpf:"ip_ifindex_map"`
}

// tcObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadTcObjects or ebpf.CollectionSpec.LoadAndAssign.
type tcObjects struct {
	tcPrograms
	tcMaps
}

func (o *tcObjects) Close() error {
	return _TcClose(
		&o.tcPrograms,
		&o.tcMaps,
	)
}

// tcMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadTcObjects or ebpf.CollectionSpec.LoadAndAssign.
type tcMaps struct {
	IpIfindexMap *ebpf.Map `ebpf:"ip_ifindex_map"`
}

func (m *tcMaps) Close() error {
	return _TcClose(
		m.IpIfindexMap,
	)
}

// tcPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadTcObjects or ebpf.CollectionSpec.LoadAndAssign.
type tcPrograms struct {
	TcIngress *ebpf.Program `ebpf:"tc_ingress"`
}

func (p *tcPrograms) Close() error {
	return _TcClose(
		p.TcIngress,
	)
}

func _TcClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed tc_bpfel.o
var _TcBytes []byte
