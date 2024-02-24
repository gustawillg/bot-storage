package main

import (
	"flag"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// Implementação do FileSystem
type FS struct{}

// Implementação do método Root
func (f FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

// Implementação do Directory
type Dir struct{}

// Implementação dos métodos do Diretório
func (d Dir) Attr() fuse.Attr {
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0755}
}

func main() {
	flag.Parse()

	// Montar o FileSystem
	c, err := fuse.Mount(
		os.Args[1],
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Servir o FileSystem
	err = fs.Serve(c, FS{})
	if err != nil {
		log.Fatal(err)
	}

	// Esperar até que seja desmontado
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}
