package mysqlmodels

import (
	"errors"
)

// Get the list of all the registered processes
func (u *DbModel) ProcessList() ([]Process, error) {
	rows, err := u.Db.Find(&Process{}).Rows()
	if err != nil {
		return nil, err
	}
	var pidList []Process
	for rows.Next() {
		var process Process
		u.Db.ScanRows(rows, &process)
		pidList = append(pidList, process)
	}

	if len(pidList) == 0 {
		return nil, errors.New("process: process list is empty")
	}
	return pidList, nil
}

// Get the information of the process id
func (u *DbModel) GetProcessInfo(id string) (*Process, error) {
	var process Process
	queryRes := u.Db.First(&process, id)
	if queryRes.Error != nil {
		return nil, queryRes.Error
	}
	return &process, nil
}
