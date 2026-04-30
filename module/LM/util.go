package LM

func (this *LockManager) MovePre(src string, dest string) error {
	_, err := this.Lock(src, 2)
	if err != nil {
		return err
	}

	_, err = this.Lock(dest, 3)
	if err != nil {
		this.Unlock(src)
		return err
	}

	return nil
}

func (this *LockManager) MovePost(src string, dest string) error {
	this.Unlock(src)
	this.Unlock(dest)
	return nil
}

func (this *LockManager) CopyPre(src string, dest string) error {
	_, err := this.Lock(src, 2)
	if err != nil {
		return err
	}

	_, err = this.Lock(dest, 3)
	if err != nil {
		this.Unlock(src)
		return err
	}

	return nil
}

func (this *LockManager) CopyPost(src string, dest string) error {
	this.Unlock(src)
	this.Unlock(dest)
	return nil
}

func (this *LockManager) RenamePre(src string, dest string) error {
	_, err := this.Lock(src, 2)
	if err != nil {
		return err
	}

	_, err = this.Lock(dest, 3)
	if err != nil {
		this.Unlock(src)
		return err
	}

	return nil
}

func (this *LockManager) RenamePost(src string, dest string) error {
	this.Unlock(src)
	this.Unlock(dest)
	return nil
}

func (this *LockManager) DeletePre(path string) error {
	_, err := this.Lock(path, 4)
	return err
}

func (this *LockManager) DeletePost(path string) error {
	return this.Unlock(path)
}

func (this *LockManager) ReadPre(path string) error {
	_, err := this.Lock(path, 1)
	return err
}

func (this *LockManager) ReadPost(path string) error {
	return this.Unlock(path)
}

func (this *LockManager) UploadPre(path string) error {
	_, err := this.Lock(path, 3)
	return err
}

func (this *LockManager) UploadPost(path string) error {
	return this.Unlock(path)
}

func (this *LockManager) DownloadPre(path string) error {
	_, err := this.Lock(path, 1)
	return err
}

func (this *LockManager) DownloadPost(path string) error {
	return this.Unlock(path)
}
