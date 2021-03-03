package mysqlmodels

import (
	"errors"
)

// Get the list of all the registered processes
func (u *DbModel) ProcessList() ([]Process, error) {
	var processList []Process
	queryResult := u.Db.Find(&processList)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}
	if len(processList) == 0 {
		return nil, errors.New("process: process list is empty")
	}
	return processList, nil
}

// Get the information of the process id
func (u *DbModel) GetProcessInfo(id int) (*Process, error) {
	var process Process
	queryRes := u.Db.First(&process, id)
	if queryRes.Error != nil {
		return nil, queryRes.Error
	}
	return &process, nil
}

func (u *DbModel) DeleteProcess(id int) error {
	process, err := u.GetProcessInfo(id)
	if err != nil {
		return err
	}
	u.Db.Delete(&process)
	return nil
}

func (u *DbModel) UpdateDescription(id int, value string) error {
	process, err := u.GetProcessInfo(id)
	if err != nil {
		return err
	}
	process.Description = value
	u.Db.Save(process)
	return nil
}

func (u *DbModel) AddProcess(name, description string) {
	process := Process{Name: name, Description: description}
	u.Db.Save(&process)
}
