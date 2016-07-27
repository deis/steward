package k8s

import (
	"fmt"

	kcl "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/runtime"
)

type errSerializerNotRegistered struct {
	contentType string
}

func (e errSerializerNotRegistered) Error() string {
	return fmt.Sprintf("serializer for %s not registered", e.contentType)
}

type errStreamingSerializerNotRegistered struct {
	contentType string
}

func (e errStreamingSerializerNotRegistered) Error() string {
	return fmt.Sprintf("streaming serializer for %s not registered", e.contentType)
}

// returns a (k8s.io/kubernetes/pkg/runtime).Serializer implementation for ServicePlanClaim. Much of this code is taken from the 'createSerializers' function in https://github.com/kubernetes/kubernetes/blob/master/pkg/client/restclient/client.go
func servicePlanClaimSerializer(config restclient.ContentConfig) (*restclient.Serializers, error) {
	negotiated := config.NegotiatedSerializer
	contentType := config.ContentType
	info, ok := negotiated.SerializerForMediaType(contentType, nil)
	if !ok {
		return nil, errSerializerNotRegistered{contentType: contentType}
	}
	streamInfo, ok := negotiated.StreamingSerializerForMediaType(contentType, nil)
	if !ok {
		return nil, errStreamingSerializerNotRegistered{contentType: contentType}
	}
	internalGV := kcl.GroupVersion{
		Group:   config.GroupVersion.Group,
		Version: runtime.APIVersionInternal,
	}
	return &restclient.Serializers{
		Encoder:             negotiated.EncoderForVersion(info.Serializer, *config.GroupVersion),
		Decoder:             negotiated.DecoderToVersion(info.Serializer, internalGV), // TODO: make this specific to ServicePlanClaim
		StreamingSerializer: streamInfo.Serializer,
		Framer:              streamInfo.Framer,
		RenegotiatedDecoder: func(contentType string, params map[string]string) (runtime.Decoder, error) {
			renegotiated, ok := negotiated.SerializerForMediaType(contentType, params)
			if !ok {
				return nil, fmt.Errorf("serializer for %s not registered", contentType)
			}
			return negotiated.DecoderToVersion(renegotiated.Serializer, internalGV), nil
		},
	}, nil
}
