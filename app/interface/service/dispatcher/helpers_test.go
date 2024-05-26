//nolint:errchkjson,goerr113 // To avoid test code in common space
package dispatcher

import (
	"encoding/json"
	"errors"
	"faraway/wow/pkg/protocol"
	"faraway/wow/pkg/test"
	"testing"

	mock "faraway/wow/app/interface/service/dispatcher/mock"

	"go.uber.org/mock/gomock"
)

type testStruct struct {
	Field string `json:"field"`
}

func Test_GetTypeName(t *testing.T) {
	t.Parallel()

	typeName := GetTypeName[testStruct]()
	expectedTypeName := "testStruct"

	test.Check(t, "Type name", expectedTypeName, typeName)
}

func Test_EncodePackage_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncoder := mock.NewMockEncoder(ctrl)
	data := &testStruct{Field: "value"}

	payloadBytes, _ := json.Marshal(data)
	mockEncoder.EXPECT().Marshal(data).Return(payloadBytes, nil)

	encodedPkg := protocol.Package{
		Tag:     "testStruct",
		Payload: payloadBytes,
	}
	encodedPkgBytes, _ := json.Marshal(encodedPkg)
	mockEncoder.EXPECT().Marshal(&encodedPkg).Return(encodedPkgBytes, nil)

	encoded, err := EncodePackage(data, mockEncoder)
	test.Nil(t, "Encode error", err)

	var pkg protocol.Package
	err = json.Unmarshal(encoded, &pkg)
	test.Nil(t, "Unmarshal package error", err)
	test.Check(t, "Package tag", "testStruct", pkg.Tag)

	var payload testStruct
	err = json.Unmarshal(pkg.Payload, &payload)
	test.Nil(t, "Unmarshal payload error", err)
	test.Check(t, "Payload", *data, payload)
}

func Test_EncodePackage_MarshalPayloadError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncoder := mock.NewMockEncoder(ctrl)
	data := &testStruct{Field: "value"}

	expectedErr := errors.New("marshal error")
	mockEncoder.EXPECT().Marshal(data).Return(nil, expectedErr)

	_, err := EncodePackage(data, mockEncoder)
	test.Err(t, "Marshal payload error", expectedErr, err)
}

func Test_EncodePackage_MarshalPackageError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncoder := mock.NewMockEncoder(ctrl)
	data := &testStruct{Field: "value"}

	payloadBytes, _ := json.Marshal(data)
	mockEncoder.EXPECT().Marshal(data).Return(payloadBytes, nil)

	expectedErr := errors.New("marshal error")
	mockEncoder.EXPECT().Marshal(gomock.Any()).Return(nil, expectedErr)

	_, err := EncodePackage(data, mockEncoder)
	test.Err(t, "Marshal package error", expectedErr, err)
}
