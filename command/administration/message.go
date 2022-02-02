package administration

import (
	"fmt"
	"strconv"
	"strings"

	proto "github.com/Jille/raftadmin/proto"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GetMessageByName(command string, args []string) (protoreflect.Message, protoreflect.MethodDescriptor, error) {

	methods := proto.File_raftadmin_proto.Services().ByName("RaftAdmin").Methods()

	// look up the user-provided command as CamelCase, kebab-case, snake_case etc.
	m := methods.ByName(protoreflect.Name(command))
	if m == nil {
		m = methods.ByName(protoreflect.Name(strcase.ToCamel(command)))
	}
	if m == nil {
		m = methods.ByName(protoreflect.Name(strcase.ToKebab(command)))
	}
	if m == nil {
		m = methods.ByName(protoreflect.Name(strcase.ToLowerCamel(command)))
	}
	if m == nil {
		m = methods.ByName(protoreflect.Name(strcase.ToSnake(command)))
	}
	if m == nil {
		return nil, nil, fmt.Errorf("unknown command %q", command)
	}

	// sort fields by field number
	reqDesc := m.Input()
	unorderedFields := reqDesc.Fields()
	fields := make([]protoreflect.FieldDescriptor, unorderedFields.Len())
	for i := 0; unorderedFields.Len() > i; i++ {
		f := unorderedFields.Get(i)
		fields[f.Number()-1] = f
	}
	if len(args) != len(fields) {
		var names []string
		for _, f := range fields {
			names = append(names, fmt.Sprintf("<%s>", f.TextName()))
		}
		return nil, nil, fmt.Errorf("invalid command arguments: '%s' requires %s", command, strings.Join(names, " "))
	}

	// convert given strings to the right type and set them on the request proto
	req := messageFromDescriptor(reqDesc)
	for i, f := range fields {
		s := args[i]
		var v protoreflect.Value
		switch f.Kind() {
		case protoreflect.StringKind:
			v = protoreflect.ValueOfString(s)
		case protoreflect.BytesKind:
			v = protoreflect.ValueOfBytes([]byte(s))
		case protoreflect.Uint64Kind:
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return nil, nil, err
			}
			v = protoreflect.ValueOfUint64(uint64(i))
		default:
			return nil, nil, fmt.Errorf("internal error: kind %s is not yet supported", f.Kind().String())
		}
		req.Set(f, v)
	}
	return req, m, nil
}

// There is no way to go from a protoreflect.MessageDescriptor to an instance of the message :(
var prototypes = []protoreflect.ProtoMessage{
	&proto.Future{},
	&proto.AwaitResponse{},
	&proto.ForgetResponse{},
	&proto.AddVoterRequest{},
	&proto.AddNonvoterRequest{},
	&proto.ApplyLogRequest{},
	&proto.AppliedIndexRequest{},
	&proto.AppliedIndexResponse{},
	&proto.BarrierRequest{},
	&proto.DemoteVoterRequest{},
	&proto.GetConfigurationRequest{},
	&proto.GetConfigurationResponse{},
	&proto.LastContactRequest{},
	&proto.LastContactResponse{},
	&proto.LastIndexRequest{},
	&proto.LastIndexResponse{},
	&proto.LeaderRequest{},
	&proto.LeaderResponse{},
	&proto.LeadershipTransferRequest{},
	&proto.LeadershipTransferToServerRequest{},
	&proto.RemoveServerRequest{},
	&proto.ShutdownRequest{},
	&proto.SnapshotRequest{},
	&proto.StateRequest{},
	&proto.StateResponse{},
	&proto.StatsRequest{},
	&proto.StatsResponse{},
	&proto.VerifyLeaderRequest{},
}

// messageFromDescriptor creates a new Message for a MessageDescriptor.
func messageFromDescriptor(d protoreflect.MessageDescriptor) protoreflect.Message {
	for _, m := range prototypes {
		if m.ProtoReflect().Descriptor() == d {
			return m.ProtoReflect().New()
		}
	}
	panic(fmt.Errorf("unknown type %q; please add it to prototypes", d.FullName()))
}
