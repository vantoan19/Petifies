// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class ImageModel {
  final String uri;
  final int width;
  final int height;
  ImageModel({
    required this.uri,
    required this.width,
    required this.height,
  });

  ImageModel copyWith({
    String? uri,
    int? width,
    int? height,
  }) {
    return ImageModel(
      uri: uri ?? this.uri,
      width: width ?? this.width,
      height: height ?? this.height,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'uri': uri,
      'width': width,
      'height': height,
    };
  }

  factory ImageModel.fromMap(Map<String, dynamic> map) {
    return ImageModel(
      uri: map['uri'] as String,
      width: map['width'] as int,
      height: map['height'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory ImageModel.fromJson(String source) =>
      ImageModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() => 'ImageModel(uri: $uri, width: $width, height: $height)';

  @override
  bool operator ==(covariant ImageModel other) {
    if (identical(this, other)) return true;

    return other.uri == uri && other.width == width && other.height == height;
  }

  @override
  int get hashCode => uri.hashCode ^ width.hashCode ^ height.hashCode;
}
