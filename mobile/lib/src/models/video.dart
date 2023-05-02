// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class NetworkVideoModel {
  final String uri;
  final String description;
  NetworkVideoModel({
    required this.uri,
    this.description = "",
  });

  NetworkVideoModel copyWith({
    String? uri,
    String? description,
  }) {
    return NetworkVideoModel(
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

  factory NetworkVideoModel.fromMap(Map<String, dynamic> map) {
    return NetworkVideoModel(
      uri: map['uri'] as String,
      description: map['description'] as String,
    );
  }

  String toJson() => json.encode(toMap());

  factory NetworkVideoModel.fromJson(String source) =>
      NetworkVideoModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() =>
      'NetworkVideoModel(uri: $uri, description: $description)';

  @override
  bool operator ==(covariant NetworkVideoModel other) {
    if (identical(this, other)) return true;

    return other.uri == uri && other.description == description;
  }

  @override
  int get hashCode => uri.hashCode ^ description.hashCode;
}
