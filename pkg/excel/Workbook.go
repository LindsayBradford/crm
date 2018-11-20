// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Workbook interface {
	Worksheets() (worksheets Worksheets)
	Worksheet(index uint) Worksheet
	WorksheetNamed(name string) Worksheet
	Save()
	SaveAs(newFileName string)
	Close(args ...interface{})
	SetProperty(propertyName string, propertyValue interface{})
	Release()
}

type WorkbookImpl struct {
	oleWrapper
}

func (wb *WorkbookImpl) WithDispatch(dispatch *ole.IDispatch) *WorkbookImpl {
	wb.dispatch = dispatch
	return wb
}

func (wb *WorkbookImpl) Worksheets() (worksheets Worksheets) {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot retrieve worksheets"))
		}
	}()

	dispatch := wb.getProperty("Worksheets")
	return new(WorksheetsImpl).WithDispatch(dispatch)
}

func (wb *WorkbookImpl) Worksheet(index uint) Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open worksheet at index [" + fmt.Sprintf("%d", index) + "]"))
		}
	}()

	dispatch := wb.getProperty("Worksheets", index)
	return new(WorksheetImpl).WithDispatch(dispatch)
}

func (wb *WorkbookImpl) WorksheetNamed(name string) Worksheet {
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New("cannot open worksheet [" + name + "]"))
		}
	}()

	dispatch := wb.getProperty("Worksheets", name)
	return new(WorksheetImpl).WithDispatch(dispatch)
}

func (wb *WorkbookImpl) Save() {
	wb.call("Save")
}

func (wb *WorkbookImpl) SaveAs(newFileName string) {
	wb.call("SaveAs", newFileName)
}

func (wb *WorkbookImpl) Close(parameters ...interface{}) {
	wb.call("Close", parameters...)
	wb.Release()
}

func (wb *WorkbookImpl) getProperty(propertyName string, parameters ...interface{}) *ole.IDispatch {
	return getProperty(wb.dispatch, propertyName, parameters...)
}

func (wb *WorkbookImpl) call(methodName string, parameters ...interface{}) *ole.IDispatch {
	return callMethod(wb.dispatch, methodName, parameters...)
}

func (wb *WorkbookImpl) SetProperty(propertyName string, propertyValue interface{}) {
	setProperty(wb.dispatch, propertyName, propertyValue)
}