// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class NetworkImageModel {
  final String uri;
  final String description;
  NetworkImageModel({
    required this.uri,
    this.description = "",
  });

  NetworkImageModel copyWith({
    String? uri,
    String? description,
  }) {
    return NetworkImageModel(
      uri: uri ?? this.uri,
      description: description ?? this.description,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'uri': uri,
      'description': description,
    };
  }

  factory NetworkImageModel.fromMap(Map<String, dynamic> map) {
    return NetworkImageModel(
      uri: map['uri'] as String,
      description: map['description'] as String,
    );
  }

  String toJson() => json.encode(toMap());

  factory NetworkImageModel.fromJson(String source) =>
      NetworkImageModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() =>
      'NetworkImageModel(uri: $uri, description: $description)';

  @override
  bool operator ==(covariant NetworkImageModel other) {
    if (identical(this, other)) return true;

    return other.uri == uri && other.description == description;
  }

  @override
  int get hashCode => uri.hashCode ^ description.hashCode;
}
