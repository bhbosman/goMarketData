// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/goMarketData/internal/googleSheets (interfaces: IGoogleSheets)

// Package googleSheets is a generated GoMock package.
package googleSheets

import (
	"context"
	fmt "fmt"

	errors "github.com/bhbosman/gocommon/errors"
)

// Interface A Comment
// Interface github.com/bhbosman/goMarketData/internal/googleSheets
// Interface IGoogleSheets
// Interface IGoogleSheets, Method: MultiSend
type IGoogleSheetsMultiSendIn struct {
	arg0 []interface{}
}

type IGoogleSheetsMultiSendOut struct {
}
type IGoogleSheetsMultiSendError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IGoogleSheetsMultiSendError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IGoogleSheetsMultiSend struct {
	inData         IGoogleSheetsMultiSendIn
	outDataChannel chan IGoogleSheetsMultiSendOut
}

func NewIGoogleSheetsMultiSend(waitToComplete bool, arg0 ...interface{}) *IGoogleSheetsMultiSend {
	var outDataChannel chan IGoogleSheetsMultiSendOut
	if waitToComplete {
		outDataChannel = make(chan IGoogleSheetsMultiSendOut)
	} else {
		outDataChannel = nil
	}
	return &IGoogleSheetsMultiSend{
		inData: IGoogleSheetsMultiSendIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IGoogleSheetsMultiSend) Wait(onError func(interfaceName string, methodName string, err error) error) (IGoogleSheetsMultiSendOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IGoogleSheetsMultiSendError{
			InterfaceName: "IGoogleSheets",
			MethodName:    "MultiSend",
			Reason:        "Channel for IGoogleSheets::MultiSend returned false",
		}
		if onError != nil {
			err := onError("IGoogleSheets", "MultiSend", generatedError)
			return IGoogleSheetsMultiSendOut{}, err
		} else {
			return IGoogleSheetsMultiSendOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IGoogleSheetsMultiSend) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIGoogleSheetsMultiSend(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 ...interface{}) (IGoogleSheetsMultiSendOut, error) {
	if context != nil && context.Err() != nil {
		return IGoogleSheetsMultiSendOut{}, context.Err()
	}
	data := NewIGoogleSheetsMultiSend(waitToComplete, arg0...)
	if waitToComplete {
		defer func(data *IGoogleSheetsMultiSend) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IGoogleSheetsMultiSendOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IGoogleSheetsMultiSendOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IGoogleSheetsMultiSendOut{}, err
	}
	return v, nil
}

// Interface IGoogleSheets, Method: Send
type IGoogleSheetsSendIn struct {
	arg0 interface{}
}

type IGoogleSheetsSendOut struct {
	Args0 error
}
type IGoogleSheetsSendError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *IGoogleSheetsSendError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type IGoogleSheetsSend struct {
	inData         IGoogleSheetsSendIn
	outDataChannel chan IGoogleSheetsSendOut
}

func NewIGoogleSheetsSend(waitToComplete bool, arg0 interface{}) *IGoogleSheetsSend {
	var outDataChannel chan IGoogleSheetsSendOut
	if waitToComplete {
		outDataChannel = make(chan IGoogleSheetsSendOut)
	} else {
		outDataChannel = nil
	}
	return &IGoogleSheetsSend{
		inData: IGoogleSheetsSendIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *IGoogleSheetsSend) Wait(onError func(interfaceName string, methodName string, err error) error) (IGoogleSheetsSendOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &IGoogleSheetsSendError{
			InterfaceName: "IGoogleSheets",
			MethodName:    "Send",
			Reason:        "Channel for IGoogleSheets::Send returned false",
		}
		if onError != nil {
			err := onError("IGoogleSheets", "Send", generatedError)
			return IGoogleSheetsSendOut{}, err
		} else {
			return IGoogleSheetsSendOut{}, generatedError
		}
	}
	return data, nil
}

func (self *IGoogleSheetsSend) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallIGoogleSheetsSend(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 interface{}) (IGoogleSheetsSendOut, error) {
	if context != nil && context.Err() != nil {
		return IGoogleSheetsSendOut{}, context.Err()
	}
	data := NewIGoogleSheetsSend(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *IGoogleSheetsSend) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return IGoogleSheetsSendOut{}, context.Err()
	}
	channel <- data
	var err error
	var v IGoogleSheetsSendOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return IGoogleSheetsSendOut{}, err
	}
	return v, nil
}

func ChannelEventsForIGoogleSheets(next IGoogleSheets, event interface{}) (bool, error) {
	switch v := event.(type) {
	case *IGoogleSheetsMultiSend:
		data := IGoogleSheetsMultiSendOut{}
		next.MultiSend(v.inData.arg0...)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	case *IGoogleSheetsSend:
		data := IGoogleSheetsSendOut{}
		data.Args0 = next.Send(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	default:
		return false, nil
	}
}
