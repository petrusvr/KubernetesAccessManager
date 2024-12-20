package cluster

import (
	"context"
	"errors"
	"testing"
	"github.com/ZPI-2024-25/KubernetesAccessManager/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

type MockResourceInterface struct {
	dynamic.ResourceInterface
	mock.Mock
	ReturnedValue *unstructured.Unstructured
	ReturnedError error
}

func (m *MockResourceInterface) Get(ctx context.Context, name string, 
	options metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedValue, m.ReturnedError
}

func (m *MockResourceInterface) Create(ctx context.Context, obj *unstructured.Unstructured, 
	options metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedValue, m.ReturnedError
}

func (m *MockResourceInterface) Delete(ctx context.Context, name string, 
	options metav1.DeleteOptions, subresources ...string) error {
	return m.ReturnedError
}

func (m *MockResourceInterface) Update(ctx context.Context, obj *unstructured.Unstructured, 
	options metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return m.ReturnedValue, m.ReturnedError
}

func MockResourceDetailsUnstructured(mockResource map[string]interface{}) *models.ResourceDetails {
	var castedResourceDetails interface{} = &unstructured.Unstructured{Object: mockResource}
	return &models.ResourceDetails{ResourceDetails: &castedResourceDetails}
}

func MockResourceDetailsMap(mockResource map[string]interface{}) *models.ResourceDetails {
	var castedResourceDetails interface{} = mockResource
	return &models.ResourceDetails{ResourceDetails: &castedResourceDetails}
}

func MockResourceDetails() *models.ResourceDetails {
	return MockResourceDetailsUnstructured(map[string]interface{}{
		"key": "value",
		"namespace": "validNamespace",
	})}

func MockUnstructured() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
		}}}


func TestGetResourceError(t *testing.T) {

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test GetResource", func(t *testing.T) {
		_, err := GetResource("Pod", "validNamespace", "validName", getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestGetResourceSuccess(t *testing.T) {

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured()}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test GetResource", func(t *testing.T) {
		result, err := GetResource("Pod", "validNamespace", "validName", getResourceI)
		assert.Nil(t, err)
		expected := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
		}
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expected {
			assert.Equal(t, value, resultObj.Object[key])
		}
	})
}

func TestGetResourceErrorFromGet(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured(), ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, 
		emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test GetResourceErrorFromGet", func(t *testing.T) {
		res, err := GetResource("Pod", "validNamespace", "validName", getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Code, 500)
		assert.EqualValues(t, err.Message, "Internal server error: error")
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}

func TestCreateResourceError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, 
		emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test CreateResourceError", func(t *testing.T) {
		_, err := CreateResource("Pod", "validNamespace", models.ResourceDetails{}, getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestCreateResourceSuccess(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured()}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test CreateResource", func(t *testing.T) {
		result, err := CreateResource("Pod", "validNamespace", *MockResourceDetailsMap(map[string]interface{}{
			"key": "value", 
			"namespace": "validNamespace",
			}), getResourceI)
		assert.Nil(t, err)

		expectedResourceDetails := MockResourceDetails()
		expectedObj := (*expectedResourceDetails.ResourceDetails).(*unstructured.Unstructured)
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expectedObj.Object {
			assert.Equal(t, value, resultObj.Object[key])
		}
	})
}

func TestCreateResourceErrorFromCreate(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured(), ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, 
		emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test CreateResourceErrorFromCreate", func(t *testing.T) {
		res, err := CreateResource("Pod", "validNamespace", *MockResourceDetailsMap(map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
		}), getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 500, err.Code)
		assert.EqualValues(t, "Internal server error: error", err.Message)
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}

func TestDeleteResourceError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test DeleteResourceError", func(t *testing.T) {
		err := DeleteResource("Pod", "validNamespace", "validName", getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestDeleteResourceSuccess(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test DeleteResource", func(t *testing.T) {
		err := DeleteResource("Pod", "validNamespace", "validName", getResourceI)
		assert.Nil(t, err)
	})
}

func TestDeleteResourceErrorFromDelete(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test DeleteResourceErrorFromDelete", func(t *testing.T) {
		err := DeleteResource("Pod", "validNamespace", "validName", getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 500, err.Code)
		assert.EqualValues(t, "Internal server error: error", err.Message)
	})
}

func TestUpdateResourceError(t *testing.T) {
	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).Return(&MockResourceInterface{}, &models.ModelError{Code: 404, Message: "Not found"})
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), args.Get(1).(*models.ModelError)
	}

	t.Run("Test UpdateResourceError", func(t *testing.T) {
		_, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetails(), getResourceI)
		assert.NotNil(t, err)
		assert.Equal(t, &models.ModelError{Code: 404, Message: "Not found"}, err)
	})
}

func TestUpdateResourceSuccess(t *testing.T) {
	expected := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
			"metadata": map[string]interface{}{"name": "validName"},
		}

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: &unstructured.Unstructured{Object: expected}}, nil)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test UpdateResource", func(t *testing.T) {
		result, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetailsMap(expected), getResourceI)
		assert.Nil(t, err)
		resultObj := (*result.ResourceDetails).(*unstructured.Unstructured)
		for key, value := range expected {
			assert.EqualValues(t, value, resultObj.Object[key])
		}
	})
}

func TestUpdateResourceErrorFromUpdate (t *testing.T) {
	dummy := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
			"metadata": map[string]interface{}{"name": "validName"},
		}

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: MockUnstructured(), ReturnedError: errors.New("error")}, nil)
	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}
	t.Run("Test UpdateResourceErrorFromUpdate", func(t *testing.T) {
		res, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetailsMap(dummy), getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 500, err.Code)
		assert.EqualValues(t, "Internal server error: error", err.Message)
		assert.EqualValues(t, res, models.ResourceDetails{})
	})
}

func TestUpdateResourceWrongName (t *testing.T) {
	expected := map[string]interface{}{
			"key": "value",
			"namespace": "validNamespace",
			"metadata": map[string]interface{}{"name": "differentNameThanExpected"},
		}

	mockGetResourceI := new(mock.Mock)
	mockGetResourceI.On("func1", "Pod", "validNamespace", DefaultNamespace).
		Return(&MockResourceInterface{ReturnedValue: &unstructured.Unstructured{Object: expected}}, nil)

	getResourceI := func(resourceType string, namespace string, emptyNamespace string) (dynamic.ResourceInterface, *models.ModelError) {
		args := mockGetResourceI.Called(resourceType, namespace, emptyNamespace)
		return args.Get(0).(dynamic.ResourceInterface), nil
	}

	t.Run("Test UpdateResource", func(t *testing.T) {
		_, err := UpdateResource("Pod", "validNamespace", "validName", *MockResourceDetailsMap(expected), getResourceI)

		assert.NotNil(t, err)
		assert.EqualValues(t, 400, err.Code)
		assert.EqualValues(t, "Invalid Input: Different resource names", err.Message)
	})
}








