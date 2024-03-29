// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package sciensano

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson320b91e2DecodeGithubComClambinSciensanoInternalSciensano(in *jlexer.Lexer, out *Vaccinations) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Vaccinations, 0, 0)
			} else {
				*out = Vaccinations{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Vaccination
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson320b91e2EncodeGithubComClambinSciensanoInternalSciensano(out *jwriter.Writer, in Vaccinations) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Vaccinations) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson320b91e2EncodeGithubComClambinSciensanoInternalSciensano(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vaccinations) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson320b91e2EncodeGithubComClambinSciensanoInternalSciensano(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vaccinations) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson320b91e2DecodeGithubComClambinSciensanoInternalSciensano(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vaccinations) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson320b91e2DecodeGithubComClambinSciensanoInternalSciensano(l, v)
}
func easyjson320b91e2DecodeGithubComClambinSciensanoInternalSciensano1(in *jlexer.Lexer, out *Vaccination) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "DATE":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.TimeStamp).UnmarshalJSON(data))
			}
		case "BRAND":
			out.Manufacturer = string(in.String())
		case "REGION":
			out.Region = string(in.String())
		case "AGEGROUP":
			out.AgeGroup = string(in.String())
		case "SEX":
			out.Gender = string(in.String())
		case "DOSE":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Dose).UnmarshalJSON(data))
			}
		case "COUNT":
			out.Count = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson320b91e2EncodeGithubComClambinSciensanoInternalSciensano1(out *jwriter.Writer, in Vaccination) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"DATE\":"
		out.RawString(prefix[1:])
		out.Raw((in.TimeStamp).MarshalJSON())
	}
	{
		const prefix string = ",\"BRAND\":"
		out.RawString(prefix)
		out.String(string(in.Manufacturer))
	}
	{
		const prefix string = ",\"REGION\":"
		out.RawString(prefix)
		out.String(string(in.Region))
	}
	{
		const prefix string = ",\"AGEGROUP\":"
		out.RawString(prefix)
		out.String(string(in.AgeGroup))
	}
	{
		const prefix string = ",\"SEX\":"
		out.RawString(prefix)
		out.String(string(in.Gender))
	}
	{
		const prefix string = ",\"DOSE\":"
		out.RawString(prefix)
		out.Raw((in.Dose).MarshalJSON())
	}
	{
		const prefix string = ",\"COUNT\":"
		out.RawString(prefix)
		out.Int(int(in.Count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Vaccination) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson320b91e2EncodeGithubComClambinSciensanoInternalSciensano1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vaccination) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson320b91e2EncodeGithubComClambinSciensanoInternalSciensano1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vaccination) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson320b91e2DecodeGithubComClambinSciensanoInternalSciensano1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vaccination) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson320b91e2DecodeGithubComClambinSciensanoInternalSciensano1(l, v)
}
