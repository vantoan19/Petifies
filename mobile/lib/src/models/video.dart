// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class VideoModel {
  final String uri;
  final int width;
  final int height;
  final int durationInSec;
  VideoModel({
    required this.uri,
    required this.width,
    required this.height,
    required this.durationInSec,
  });

  VideoModel copyWith({
    String? uri,
    int? width,
    int? height,
    int? durationInSec,
  }) {
    return VideoModel(
      uri: uri ?? this.uri,
      width: width ?? this.width,
      height: height ?? this.height,
      durationInSec: durationInSec ?? this.durationInSec,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'uri': uri,
      'width': width,
      'height': height,
      'durationInSec': durationInSec,
    };
  }

  factory VideoModel.fromMap(Map<String, dynamic> map) {
    return VideoModel(
      uri: map['uri'] as String,
      width: map['width'] as int,
      height: map['height'] as int,
      durationInSec: map['durationInSec'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory VideoModel.fromJson(String source) =>
      VideoModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'VideoModel(uri: $uri, width: $width, height: $height, durationInSec: $durationInSec)';
  }

  @override
  bool operator ==(covariant VideoModel other) {
    if (identical(this, other)) return true;

    return other.uri == uri &&
        other.width == width &&
        other.height == height &&
        other.durationInSec == durationInSec;
  }

  @override
  int get hashCode {
    return uri.hashCode ^
        width.hashCode ^
        height.hashCode ^
        durationInSec.hashCode;
  }
}
