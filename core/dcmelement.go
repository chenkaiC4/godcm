package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

// DcmElement indentified the data element tag.
type DcmElement struct {
	Tag          DcmTag
	Name         string
	VR           string
	Length       int64
	Value        []byte
	Squence      *DcmSQElement
	isExplicitVR bool
	byteOrder    EByteOrder
	isReadValue  bool
	isReadPixel  bool
}

// GetValueString convert value to string according to VR
func (e DcmElement) GetValueString() string {
	buf := bytes.NewBuffer(e.Value)
	var result string
	if e.Tag.Element == 0x0000 {
		var i int32
		binary.Read(buf, binary.LittleEndian, &i)
		result = fmt.Sprintf("%d", i)
		return result
	}
	switch e.VR {
	case "FL", "FD", "OD", "OF":
		var f float64
		if e.byteOrder == EBOBigEndian {
			binary.Read(buf, binary.BigEndian, &f)
		} else {
			binary.Read(buf, binary.LittleEndian, &f)
		}
		result = fmt.Sprintf("%f", f)
	case "OL", "SL", "SS", "UL":
		var i int32
		if e.byteOrder == EBOBigEndian {
			binary.Read(buf, binary.BigEndian, &i)
		} else {
			binary.Read(buf, binary.LittleEndian, &i)
		}
		result = fmt.Sprintf("%d", i)
	case "US", "US or SS":
		var i uint16
		if e.byteOrder == EBOBigEndian {
			binary.Read(buf, binary.BigEndian, &i)
		} else {
			binary.Read(buf, binary.LittleEndian, &i)
		}
		result = fmt.Sprintf("%d", i)
	case "AE", "AS", "CS", "DA", "DS", "DT", "IS", "LO", "LT", "PN", "ST", "UI", "UT", "TM", "SH":
		result = string(bytes.Trim(e.Value, "\x00"))
	default:
		result = fmt.Sprintf("%x", e.Value)
	}
	return strings.TrimSpace(result)
}

// String convert to string value
func (e DcmElement) String() string {
	if e.Squence != nil {
		return fmt.Sprintf("Tag:%s; VR:%s; Length:%d; Value:%s; Sequence : %v", e.Tag, e.VR, e.Length, e.GetValueString(), e.Squence)
	}
	return fmt.Sprintf("Tag:%s; VR:%s; Length:%d; Value:%s", e.Tag, e.VR, e.Length, e.GetValueString())
}

// ReadUINT16 is to read a uint16 value from the file.
func (e DcmElement) ReadUINT16(s *DcmFileStream) (uint16, error) {
	v, err := s.Read(2)
	if err != nil {
		return 0, err
	}
	var result uint16
	buf := bytes.NewReader(v)
	if e.byteOrder == EBOBigEndian {
		err = binary.Read(buf, binary.BigEndian, &result)
	} else {
		err = binary.Read(buf, binary.LittleEndian, &result)
	}

	return result, err
}

// ReadUINT32 is to read a uint32 value from the file.
func (e DcmElement) ReadUINT32(s *DcmFileStream) (uint32, error) {
	v, err := s.Read(4)
	if err != nil {
		return 0, err
	}
	var result uint32
	buf := bytes.NewReader(v)
	if e.byteOrder == EBOBigEndian {
		err = binary.Read(buf, binary.BigEndian, &result)
	} else {
		err = binary.Read(buf, binary.LittleEndian, &result)
	}
	return result, err
}

// ReadDcmTagGroup read tag  group of the dicom element.
func (e *DcmElement) ReadDcmTagGroup(s *DcmFileStream) error {
	var err error
	e.Tag.Group, err = e.ReadUINT16(s)
	if err != nil {
		return err
	}
	return nil
}

// ReadDcmTagElement read tag  group of the dicom element.
func (e *DcmElement) ReadDcmTagElement(s *DcmFileStream) error {
	var err error
	e.Tag.Element, err = e.ReadUINT16(s)
	if err != nil {
		return err
	}
	return nil
}

// ReadDcmTag is to read group and element
func (e *DcmElement) ReadDcmTag(s *DcmFileStream) error {
	err := e.ReadDcmTagGroup(s)
	if err != nil {
		return err
	}
	err = e.ReadDcmTagElement(s)
	if err != nil {
		return err
	}
	return nil
}

// ReadDcmVR is to read vr
func (e *DcmElement) ReadDcmVR(s *DcmFileStream) error {
	var err error
	e.VR, err = s.ReadString(2)
	return err
}

// ReadValueLengthWithExplicitVR gets the value length of the dicom element with explicit VR.
func (e *DcmElement) ReadValueLengthWithExplicitVR(s *DcmFileStream) error {
	switch e.VR {
	case "OB", "OD", "OF", "OL", "OW", "SQ", "UC", "UR", "UT", "UN":
		// skip the reserved 2 bytes
		_, err := s.Skip(2)
		if err != nil {
			return err
		}
		err = e.ReadValueLengthUint32(s)
		if err != nil {
			return err
		}
	default:
		// read value length
		err := e.ReadValueLengthUint16(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadValueLengthWithImplicitVR gets the value length of the dicom element with implicit VR.
func (e *DcmElement) ReadValueLengthWithImplicitVR(s *DcmFileStream) error {
	return e.ReadValueLengthUint32(s)
}

// ReadValueLengthUint16 read 2 bytes value length
func (e *DcmElement) ReadValueLengthUint16(s *DcmFileStream) error {
	l, err := e.ReadUINT16(s)
	if err != nil {
		return err
	}
	e.Length = int64(l)
	return nil
}

// ReadValueLengthUint32 read 4 bytes value length
func (e *DcmElement) ReadValueLengthUint32(s *DcmFileStream) error {
	l, err := e.ReadUINT32(s)
	if err != nil {
		return err
	}
	e.Length = int64(l)
	return nil

}

// ReadValue get or skip the element value.
func (e *DcmElement) ReadValue(s *DcmFileStream) error {
	var err error
	if !e.isReadPixel {
		if e.Tag.Group == 0x7fe0 {
			_, err = s.Skip(e.Length)
			return err
		}
	}

	if e.isReadValue {
		// read element value
		e.Value, err = s.Read(e.Length)
	} else {
		_, err = s.Skip(e.Length)
	}
	return err
}

// ReadDcmElement read one dicom element.
func (e *DcmElement) ReadDcmElement(s *DcmFileStream) error {

	if e.isExplicitVR {
		return e.ReadDcmElementWithExplicitVR(s)
	}
	return e.ReadDcmElementWithImplicitVR(s)
}

func (e *DcmElement) readDcmSQElement(s *DcmFileStream) error {
	e.Squence = new(DcmSQElement)
	err := e.Squence.Read(s, e.Length, e.isExplicitVR, e.isReadValue)
	if err != nil {
		return err
	}
	//	log.Println(e.String())
	return nil
}

// ReadDcmElementWithExplicitVR read the data element with explicit VR.
func (e *DcmElement) ReadDcmElementWithExplicitVR(s *DcmFileStream) error {
	// read dicom tag
	err := e.ReadDcmTag(s)
	if err != nil {
		return err
	}

	// read VR
	err = e.ReadDcmVR(s)
	if err != nil {
		return err
	}

	//read the value length
	err = e.ReadValueLengthWithExplicitVR(s)
	if err != nil {
		return err
	}
	// skip reading value if length is zero
	if e.Length == 0 {
		//		log.Println(e.String())
		return nil
	}

	// read sequence items
	if e.VR == "SQ" {
		return e.readDcmSQElement(s)
	}

	// read VR:UN with unknown length
	if e.VR == "UN" && e.Length == 0xFFFFFFFF {
		return e.readDcmSQElement(s)
	}

	// encapsulated pixel data
	if e.Tag == DCMPixelData && e.Length == 0xFFFFFFFF {
		return e.readDcmSQElement(s)
	}

	err = e.ReadValue(s)
	if err != nil {
		return err
	}

	//	log.Println(e.String())

	return nil
}

// ReadDcmElementWithImplicitVR read the data element with implicit VR.
func (e *DcmElement) ReadDcmElementWithImplicitVR(s *DcmFileStream) error {
	// read dciom tag
	err := e.ReadDcmTag(s)
	if err != nil {
		return err
	}

	// get VR from Dicom Element registry
	err = FindDcmElmentByTag(e)
	if err != nil {
		log.Println(err.Error())
	}

	// read the value length
	err = e.ReadValueLengthWithImplicitVR(s)
	if err != nil {
		return err
	}

	err = e.ReadValue(s)
	if err != nil {
		return err
	}

	//	log.Println(e)

	/*
		if elem.Tag.Group != 0x7fe0 {
					log.Println(elem)
		}
	*/
	return nil

}
