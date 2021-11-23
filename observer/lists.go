package observer

type ListWatcher interface {
	ChangeName(name string)
	AddFolder(path string)
	RemoveFolder(path string)
}

type List struct {
	name    string
	folders []string
}

func (ll *List) ChangeName(name string) {
	ll.name = name
}

func (ll *List) AddFolder(path string) {
	ll.folders = append(ll.folders, path)
}

func (ll *List) RemoveFolder(path string) {
	length := len(ll.folders)

	for i, folder := range ll.folders {
		if folder == path {
			ll.folders[i] = ll.folders[length-1]
			ll.folders = ll.folders[:length-1]
			break
		}
	}
}
