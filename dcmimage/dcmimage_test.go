package dcmimage

import (
	"os"
	"strconv"
	"testing"

	"github.com/grayzone/godcm/core"
	"github.com/grayzone/godcm/util"
)

type testCase struct {
	in   string
	want bool
}

func gettestdataset() []testCase {
	result := []testCase{
		{"MR-MONO2-8-16x-heart.dcm", false},
		{"US-MONO2-8-8x-execho.dcm", false},
		{"xr_tspine.dcm", false},
		{"xr_chest.dcm", false},
		{"IM0.dcm", false},
		{"image_09-12-2013_4.dcm", false},
		{"CT-MONO2-16-ankle", false},
		{"xr_chicken2.dcm", true},
		{"T23/IM-0001-0001.dcm", true},
		{"IM-0001-0010.dcm", true},
		{"GH195.dcm", true},
		{"GH064.dcm", true},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_undefined_length.dcm", false},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_ELE_undefinded_length.dcm", false},
		{"US-RGB-8-esopecho.dcm", false},
	}
	return result
}

func readpixel(t *testing.T, filename string, want bool) DcmImage {
	var reader core.DcmReader
	reader.IsReadPixel = true
	reader.IsReadValue = true
	filepath := util.GetTestDataFolder() + filename
	err := reader.ReadFile(filepath)

	isCompressed, err := reader.IsCompressed()
	if err != nil {
		t.Errorf("readpixel() %s", err.Error())
	}

	isBigEndian, err := reader.IsBigEndian()
	if err != nil {
		t.Errorf("readpixel() %s", err.Error())
	}

	pixeldata := reader.Dataset.PixelData()

	var img DcmImage

	img.IsCompressed = isCompressed
	if want != img.IsCompressed {
		t.Errorf("readpixel(%s), isCompressed want %v got %v", filename, want, isCompressed)
	}

	img.IsBigEndian = isBigEndian

	var num interface{}

	num, _ = strconv.ParseUint(reader.Dataset.BitsAllocated(), 10, 16)
	img.BitsAllocated = uint16(num.(uint64))

	num, _ = strconv.ParseUint(reader.Dataset.Columns(), 10, 32)
	img.Columns = uint32(num.(uint64))

	num, _ = strconv.ParseUint(reader.Dataset.Rows(), 10, 32)
	img.Rows = uint32(num.(uint64))

	num, _ = strconv.ParseUint(reader.Dataset.HighBit(), 10, 16)
	img.HighBit = uint16(num.(uint64))

	num, _ = strconv.ParseFloat(reader.Dataset.WindowCenter(), 64)
	img.WindowCenter = num.(float64)

	num, _ = strconv.ParseFloat(reader.Dataset.WindowWidth(), 64)
	img.WindowWidth = num.(float64)

	num, _ = strconv.ParseFloat(reader.Dataset.RescaleIntercept(), 64)
	img.RescaleIntercept = num.(float64)

	num, _ = strconv.ParseFloat(reader.Dataset.RescaleSlope(), 64)
	img.RescaleSlope = num.(float64)

	num, _ = strconv.ParseUint(reader.Dataset.PixelRepresentation(), 10, 16)
	img.PixelRepresentation = uint16(num.(uint64))

	img.PhotometricInterpretation = reader.Dataset.PhotometricInterpretation()

	num, _ = strconv.ParseUint(reader.Dataset.NumberOfFrames(), 10, 64)
	img.NumberOfFrames = int(num.(uint64))

	num, _ = strconv.ParseUint(reader.Dataset.SamplesPerPixel(), 10, 16)
	img.SamplesPerPixel = uint16(num.(uint64))

	img.PixelData = pixeldata

	return img
}

func convert2BMP(t *testing.T, bits uint16) {
	cases := gettestdataset()
	for _, c := range cases {
		img := readpixel(t, c.in, c.want)

		for i := 0; i < img.NumberOfFrames; i++ {
			newfile := c.in + "_" + strconv.Itoa(int(bits)) + "_" + strconv.FormatUint(uint64(i), 10) + ".bmp"
			err := img.WriteBMP(newfile, bits, i)

			if err != nil {
				//		t.Errorf("ConvertToPNG() %s", err.Error())
			}

			defer os.Remove(newfile)
		}
	}
}

func TestWrite8BMP(t *testing.T) {
	convert2BMP(t, 8)
}

func TestWrite24BMP(t *testing.T) {
	convert2BMP(t, 24)
}

func TestWrite32BMP(t *testing.T) {
	convert2BMP(t, 32)
}

func TestWrite16BMP(t *testing.T) {
	convert2BMP(t, 16)
}

func TestWritePNG(t *testing.T) {
	cases := gettestdataset()
	for _, c := range cases {
		img := readpixel(t, c.in, c.want)
		for i := 0; i < img.NumberOfFrames; i++ {
			newfile := c.in + "_" + strconv.FormatUint(uint64(i), 10) + ".png"
			err := img.ConvertToPNG(newfile, i)

			if err != nil {
				//		t.Errorf("ConvertToPNG() %s", err.Error())
			}

			defer os.Remove(newfile)
		}
	}
}

func TestWriteJPG(t *testing.T) {
	cases := gettestdataset()
	for _, c := range cases {
		img := readpixel(t, c.in, c.want)

		for i := 0; i < img.NumberOfFrames; i++ {
			newfile := c.in + "_" + strconv.FormatUint(uint64(i), 10) + ".jpg"
			err := img.ConvertToJPG(newfile, i)

			if err != nil {
				//			t.Errorf("ConvertToJPG() %s", err.Error())
			}

			defer os.Remove(newfile)
		}
	}
}

/*

func TestWritePNG16(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"MR-MONO2-8-16x-heart.dcm", false},
		{"US-MONO2-8-8x-execho.dcm", false},
		{"xr_tspine.dcm", false},
		{"xr_chest.dcm", false},
		{"IM0.dcm", false},
		{"image_09-12-2013_4.dcm", false},
		{"CT-MONO2-16-ankle", false},
		{"xr_chicken2.dcm", true},
		{"T23/IM-0001-0001.dcm", true},
		{"IM-0001-0010.dcm", true},
		{"GH195.dcm", true},
		{"GH064.dcm", true},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_undefined_length.dcm", false},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_ELE_undefinded_length.dcm", false},
	}
	for _, c := range cases {
		img := readpixel(t, c.in, c.want)
		pngfile := c.in + "_16.png"
		err := img.convertToPNG16(pngfile)
		//		defer os.Remove(pngfile)
		if err != nil {
			//		t.Errorf("ConvertToPNG() %s", err.Error())
		}
	}
}


func TestWriteJPG16(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"MR-MONO2-8-16x-heart.dcm", false},
		{"US-MONO2-8-8x-execho.dcm", false},
		{"xr_tspine.dcm", false},
		{"xr_chest.dcm", false},
		{"IM0.dcm", false},
		{"image_09-12-2013_4.dcm", false},
		{"CT-MONO2-16-ankle", false},
		{"xr_chicken2.dcm", true},
		{"T23/IM-0001-0001.dcm", true},
		{"IM-0001-0010.dcm", true},
		{"GH195.dcm", true},
		{"GH064.dcm", true},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_undefined_length.dcm", false},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_ELE_undefinded_length.dcm", false},
	}
	for _, c := range cases {
		img := readpixel(t, c.in, c.want)
		jpgfile := c.in + "_16.jpg"
		err := img.convertToJPG16(jpgfile)
		defer os.Remove(jpgfile)
		if err != nil {
			//		t.Errorf("ConvertToJPG() %s", err.Error())
		}
	}
}
*/
