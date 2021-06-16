package updater

// get the currently running executable name
//

//func (u *Updater) walk(src string) error {
//	err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
//		if entry.IsDir() {
//			return nil
//		}
//
//		relativePath := strings.ReplaceAll(path, src+string(os.PathSeparator), "")
//
//		//// Check for executable and update
//		//if relativePath == u.RemoteExecutablePath {
//		//	err := u.updateExecutable(path)
//		//	if err != nil {
//		//		return err
//		//	}
//		//}
//
//		// Check for files and folders and update
//		err = u.updateFiles(src, relativePath, entry)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		return err
//	}
//	return nil
//}

//

//
//func tet() {
//	f := patcher.Folder{
//		SourcePath:      "/var/folders/hk/yl4d9gcn5v542yksl7hqk4lh0000gn/T/verbis-updater564261290/verbis/build/admin/index.html",
//		DestinationPath: "/Users/ainsley/Desktop/Reddico/tools/updater/cmd/index.html",
//		BackupPath:      "/Users/ainsley/Desktop/Reddico/tools/updater/cmd/index.html.old",
//		Matches:         []patcher.Match{},
//	}
//
//	f2 := patcher.Folder{
//		SourcePath:      "/var/folders/hk/yl4d9gcn5v542yksl7hqk4lh0000gn/T/verbis-updater564261290/verbis/build/admin",
//		DestinationPath: "/Users/ainsley/Desktop/Reddico/tools/updater/cmd/admin",
//		BackupPath:      "/Users/ainsley/Desktop/Reddico/tools/updater/cmd/admin.old",
//		Matches:         []patcher.Match{
//			{
//				"/var/folders/hk/yl4d9gcn5v542yksl7hqk4lh0000gn/T/verbis-updater728016125/verbis/build/admin/js/chunk-7efb1142.f9dddefe.js",
//				"/Users/ainsley/Desktop/Reddico/tools/updater/cmd/admin/js/chunk-7efb1142.f9dddefe.js",
//				066,
//			},
//		},
//	}
//
//	fmt.Println(f, f2)
//}
//
//
//type Patcher struct {
//	SourcePath      string
//	DestinationPath string
//	BackupPath      string
//	Mode            os.FileMode
//}
//
//
//
//type Match struct {
//	SourcePath string
//	DestinationPath string
//	Mode os.FileMode
//}
//
//
//func (f *Folder) Apply() {
//
//}
//
//
//
//func (p *Patcher) Apply() error {
//	content, err := ioutil.ReadFile(p.SourcePath)
//	if err != nil {
//		return err
//	}
//
//	//err = p.backup()
//	//if err != nil {
//	//	return err
//	//}
//
//
//	err = ioutil.WriteFile(p.DestinationPath, content, p.Mode)
//	if err != nil {
//		return p.Rollback()
//	}
//
//	return nil
//}
//
//// backup renames a directory or file to the new path.
//func (p *Patcher) backup() error {
//	// TODO check if the folder exists
//	return os.Rename(p.DestinationPath, p.BackupPath)
//}
//
//func (p *Patcher) Rollback() error {
//	return os.Rename(p.BackupPath, p.DestinationPath)
//}

//type Patcher struct {
//
//	BackupPath string
//	files []File
//}
//
//type File struct {
//	SourcePath string
//	DestinationPath string
//	Mode os.FileMode
//}
//
//func (p *Patcher) IsDirectory() bool {
//	return len(p.files) == 1 && !p.files[0].Mode.IsDir()
//}
//
//func (p *Patcher) Apply() {
//
//}
//
//
//
//
//
////// backup renames a directory or file to the new path.
////func (p *Patcher) backup() error {
////	// TODO check if the folder exists
////	return os.Rename(p.DestinationPath, p.BackupPath)
////}
////
////func (p *Patcher) Rollback() error {
////	return os.Rename(p.BackupPath, p.DestinationPath)
////}
