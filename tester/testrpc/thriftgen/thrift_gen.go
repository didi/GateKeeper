// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package thriftgen

import (
	"bytes"
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"reflect"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = reflect.DeepEqual
var _ = bytes.Equal

//Data Attributes:
//  - Text
type Data struct {
	Text string `thrift:"text,1" db:"text" json:"text"`
}

//NewData NewData
func NewData() *Data {
	return &Data{}
}

//GetText GetText
func (p *Data) GetText() string {
	return p.Text
}

//Read Read
func (p *Data) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeID, fieldID, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldID), err)
		}
		if fieldTypeID == thrift.STOP {
			break
		}
		switch fieldID {
		case 1:
			if fieldTypeID == thrift.STRING {
				if err := p.ReadField1(iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(fieldTypeID); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(fieldTypeID); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

//ReadField1 ReadField1
func (p *Data) ReadField1(iprot thrift.TProtocol) error {
	v, err := iprot.ReadString()
	if err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	}
	p.Text = v
	return nil
}

//Write Write
func (p *Data) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Data"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Data) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("text", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:text: ", p), err)
	}
	if err := oprot.WriteString(string(p.Text)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.text (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:text: ", p), err)
	}
	return err
}

//String String
func (p *Data) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Data(%+v)", *p)
}

//FormatData FormatData
type FormatData interface {
	// Parameters:
	//  - Data
	DoFormat(ctx context.Context, data *Data) (r *Data, err error)
}

//FormatDataClient FormatDataClient
type FormatDataClient struct {
	c thrift.TClient
}

//NewFormatDataClientFactory Deprecated: Use NewFormatData instead
func NewFormatDataClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *FormatDataClient {
	return &FormatDataClient{
		c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
	}
}

//NewFormatDataClientProtocol Deprecated: Use NewFormatData instead
func NewFormatDataClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *FormatDataClient {
	return &FormatDataClient{
		c: thrift.NewTStandardClient(iprot, oprot),
	}
}

//NewFormatDataClient NewFormatDataClient
func NewFormatDataClient(c thrift.TClient) *FormatDataClient {
	return &FormatDataClient{
		c: c,
	}
}

//DoFormat Parameters:
//  - Data
func (p *FormatDataClient) DoFormat(ctx context.Context, data *Data) (r *Data, err error) {
	var _args0 FormatDataDoFormatArgs
	_args0.Data = data
	var _result1 FormatDataDoFormatResult
	if err = p.c.Call(ctx, "do_format", &_args0, &_result1); err != nil {
		return
	}
	return _result1.GetSuccess(), nil
}

//FormatDataProcessor FormatDataProcessor
type FormatDataProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      FormatData
}

//AddToProcessorMap AddToProcessorMap
func (p *FormatDataProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

//GetProcessorFunction GetProcessorFunction
func (p *FormatDataProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

//ProcessorMap ProcessorMap
func (p *FormatDataProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

//NewFormatDataProcessor NewFormatDataProcessor
func NewFormatDataProcessor(handler FormatData) *FormatDataProcessor {

	self2 := &FormatDataProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self2.processorMap["do_format"] = &formatDataProcessorDoFormat{handler: handler}
	return self2
}

//Process Process
func (p *FormatDataProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqID, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqID, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqID)
	x3.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush(ctx)
	return false, x3

}

type formatDataProcessorDoFormat struct {
	handler FormatData
}

//Process Process
func (p *formatDataProcessorDoFormat) Process(ctx context.Context, seqID int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := FormatDataDoFormatArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("do_format", thrift.EXCEPTION, seqID)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush(ctx)
		return false, err
	}

	iprot.ReadMessageEnd()
	result := FormatDataDoFormatResult{}
	var retval *Data
	var err2 error
	if retval, err2 = p.handler.DoFormat(ctx, args.Data); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing do_format: "+err2.Error())
		oprot.WriteMessageBegin("do_format", thrift.EXCEPTION, seqID)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush(ctx)
		return true, err2
	}
	result.Success = retval
	if err2 = oprot.WriteMessageBegin("do_format", thrift.REPLY, seqID); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(ctx); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

//FormatDataDoFormatArgs Attributes:
//  - Data
type FormatDataDoFormatArgs struct {
	Data *Data `thrift:"data,1" db:"data" json:"data"`
}

//NewFormatDataDoFormatArgs NewFormatDataDoFormatArgs
func NewFormatDataDoFormatArgs() *FormatDataDoFormatArgs {
	return &FormatDataDoFormatArgs{}
}

//FormatDataDoFormatArgsDataDEFAULT FormatDataDoFormatArgsDataDEFAULT
var FormatDataDoFormatArgsDataDEFAULT *Data

//GetData GetData
func (p *FormatDataDoFormatArgs) GetData() *Data {
	if !p.IsSetData() {
		return FormatDataDoFormatArgsDataDEFAULT
	}
	return p.Data
}

//IsSetData IsSetData
func (p *FormatDataDoFormatArgs) IsSetData() bool {
	return p.Data != nil
}

//Read Read
func (p *FormatDataDoFormatArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeID, fieldID, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldID), err)
		}
		if fieldTypeID == thrift.STOP {
			break
		}
		switch fieldID {
		case 1:
			if fieldTypeID == thrift.STRUCT {
				if err := p.ReadField1(iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(fieldTypeID); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(fieldTypeID); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

//ReadField1 ReadField1
func (p *FormatDataDoFormatArgs) ReadField1(iprot thrift.TProtocol) error {
	p.Data = &Data{}
	if err := p.Data.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Data), err)
	}
	return nil
}

//Write Write
func (p *FormatDataDoFormatArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("do_format_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *FormatDataDoFormatArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("data", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:data: ", p), err)
	}
	if err := p.Data.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Data), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:data: ", p), err)
	}
	return err
}

//String String
func (p *FormatDataDoFormatArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("FormatDataDoFormatArgs(%+v)", *p)
}

//FormatDataDoFormatResult Attributes:
//  - Success
type FormatDataDoFormatResult struct {
	Success *Data `thrift:"success,0" db:"success" json:"success,omitempty"`
}

//NewFormatDataDoFormatResult NewFormatDataDoFormatResult
func NewFormatDataDoFormatResult() *FormatDataDoFormatResult {
	return &FormatDataDoFormatResult{}
}

//FormatDataDoFormatResultSuccessDEFAULT FormatDataDoFormatResultSuccessDEFAULT
var FormatDataDoFormatResultSuccessDEFAULT *Data

//GetSuccess GetSuccess
func (p *FormatDataDoFormatResult) GetSuccess() *Data {
	if !p.IsSetSuccess() {
		return FormatDataDoFormatResultSuccessDEFAULT
	}
	return p.Success
}

//IsSetSuccess IsSetSuccess
func (p *FormatDataDoFormatResult) IsSetSuccess() bool {
	return p.Success != nil
}

//Read Read
func (p *FormatDataDoFormatResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeID, fieldID, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldID), err)
		}
		if fieldTypeID == thrift.STOP {
			break
		}
		switch fieldID {
		case 0:
			if fieldTypeID == thrift.STRUCT {
				if err := p.ReadField0(iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(fieldTypeID); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(fieldTypeID); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

//ReadField0 ReadField0
func (p *FormatDataDoFormatResult) ReadField0(iprot thrift.TProtocol) error {
	p.Success = &Data{}
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

//Write Write
func (p *FormatDataDoFormatResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("do_format_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *FormatDataDoFormatResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

//String String
func (p *FormatDataDoFormatResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("FormatDataDoFormatResult(%+v)", *p)
}
