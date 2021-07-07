package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"time"
)

// 实例化Scheme资源注册表及Codecs编解码器
var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)

// inMediaType定义了编码类型( 即Protobuf 格式)，outMediaType 定义了解码类型(即JSON格式)。
var inMediaType = "application/vnd.kubernetes.protobuf"
var outMediaType = "application/json"

// 通过init函数将corev1资源组下的资源注册至Scheme 资源注册表中，这是因为要对Pod资源数据进行解码操作。
func init() {
	v1.AddToScheme(Scheme)
}

func main() {
	// 通过cientv3.New函数实例化Etcd Cien对象，并设置些参数
	// 例如，将Endpoints参数连接至Etcd集群的地址，将DailTimeout参数连接至集群的超时时间等。
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:4012"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 通过cli.Get函数获取Etcd集群/registry/pods/default/centos-59db99c6bc-m6tsz下的Pod资源对象数据。
	resp, err := cli.Get(context.Background(), "/registry/pods/default/centos-59db99c6bc-m6tsz")
	if err != nil {
		panic(err)
	}
	kv := resp.Kvs[0]

	// 通过newCodee函数实例化runtime.Codec编解码器，分别实例化inCodec编码器对象、outCodec解码器对象。
	inCodec := newCodec(inMediaType)
	outCodec := newCodec(outMediaType)

	// 通过runtime.Decode解码器(即protobufSerializer)解码资源对象数据并通过fmt.Println函数输出。
	obj, err := runtime.Decode(inCodec, kv.Value)
	if err != nil {
		panic(err)
	}
	fmt.Println("Decode ---")
	fmt.Println(obj)

	// 通过runtime.Encode编码器(即jsonSerializer)解码资源对象数据并通过fmt.Println函数输出。
	encoded, err := runtime.Encode(outCodec, obj)
	if err != nil {
		panic(err)
	}
	fmt.Println("Encode --- ")
	fmt.Println(string(encoded))

	/*提示: Kubernetes资源对象以二进制格式存储在Etcd 集群中，所以需要额外的解码步骤，直接从Eted 集群中获取对象的体验并不友好。
	可以通过第三方工具github.com/jpbetz/auger直接访问数据对象存储，auger用于对Kubernetes资源对象存储在Eted集群中的二进制数据进行编码和解码，
	支持将数据转换为YAML、JSON和Protobuf格式。*/
}

func newCodec(mediaTypes string) runtime.Codec {
	info, ok := runtime.SerializerInfoForMediaType(Codecs.SupportedMediaTypes(), mediaTypes)
	if !ok {
		panic(fmt.Errorf("no serializers registered for %v", mediaTypes))
	}
	cfactory := serializer.CodecFactory{}

	gv, err := schema.ParseGroupVersion("v1")
	if err != nil {
		panic("unexpected error")
	}
	encoder := cfactory.EncoderForVersion(info.Serializer, gv)
	decoder := cfactory.DecoderToVersion(info.Serializer, gv)
	return cfactory.CodecForVersions(encoder, decoder, gv, gv)

}
