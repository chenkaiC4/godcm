package dcmdata

import (
	"github.com/grayzone/godcm/ofstd"
)

/** a class handling the DICOM dataset format (files without meta header)
 */
type DcmDataset struct {
	DcmItem

	OriginalXfer E_TransferSyntax /// original transfer syntax of the dataset

	CurrentXfer E_TransferSyntax /// current transfer syntax of the dataset

}

func NewDcmDataset() *DcmDataset {
	var result DcmDataset
	result.DcmItem = *NewDcmItem(DCM_ItemTag, DCM_UndefinedLength)
	result.OriginalXfer = EXS_Unknown

	// the default transfer syntax is explicit VR with local endianness
	if GLocalByteOrder == EBO_BigEndian {
		result.CurrentXfer = EXS_BigEndianExplicit
	} else {
		result.CurrentXfer = EXS_LittleEndianExplicit
	}
	return &result
}

/** load object from a DICOM file.
 *  This method only supports DICOM objects stored as a dataset, i.e. without meta header.
 *  Use DcmFileFormat::loadFile() to load files with meta header.
 *  @param fileName name of the file to load (may contain wide chars if support enabled).
 *    Since there are various constructors for the OFFilename class, a "char *", "OFString"
 *    or "wchar_t *" can also be passed directly to this parameter.
 *  @param readXfer transfer syntax used to read the data (auto detection if EXS_Unknown)
 *  @param groupLength flag, specifying how to handle the group length tags
 *  @param maxReadLength maximum number of bytes to be read for an element value.
 *    Element values with a larger size are not loaded until their value is retrieved
 *    (with getXXX()) or loadAllDataIntoMemory() is called.
 *  @return status, EC_Normal if successful, an error code otherwise
 */
func (ds *DcmDataset) LoadFile(filename string, readXfer E_TransferSyntax, groupLength E_GrpLenEncoding, maxReadLength uint32) ofstd.OFCondition {
	err := ofstd.EC_InvalidFilename
	if len(filename) == 0 {
		return err
	}
	fileStream := NewDcmInputFileStream(filename, 0)

	err = fileStream.Status()
	if !err.Good() {
		return err
	}
	err = ds.Clear()

	if !err.Good() {
		return err
	}
	ds.TransferInit()
	err = ds.Read(fileStream, readXfer, groupLength, maxReadLength)
	ds.TransferEnd()

	//	err = fileStream

	return err
}

func (ds *DcmDataset) Read(instream *DcmInputFileStream, xfer E_TransferSyntax, glenc E_GrpLenEncoding, maxReadLength uint32) ofstd.OFCondition {
	ds.errorFlag = instream.Status()

	if ds.errorFlag.Good() && instream.Eos() {
		ds.errorFlag = EC_EndOfStream
	} else if ds.errorFlag.Good() && ds.GetTransferState() != ERW_ready {

	}

	if ds.errorFlag.Good() || ds.errorFlag == EC_EndOfStream {
		ds.errorFlag = ofstd.EC_Normal

		ds.SetTransferState(ERW_ready)
	}

	return ds.errorFlag

}
