// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class BasicUserInfoModel {
  final String id;
  final String email;
  final String? userAvatar;
  final String firstName;
  final String lastName;
  BasicUserInfoModel({
    required this.id,
    required this.email,
    this.userAvatar,
    required this.firstName,
    required this.lastName,
  });

  BasicUserInfoModel copyWith({
    String? id,
    String? email,
    String? userAvatar,
    String? firstName,
    String? lastName,
  }) {
    return BasicUserInfoModel(
      id: id ?? this.id,
      email: email ?? this.email,
      userAvatar: userAvatar ?? this.userAvatar,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'email': email,
      'userAvatar': userAvatar,
      'firstName': firstName,
      'lastName': lastName,
    };
  }

  factory BasicUserInfoModel.fromMap(Map<String, dynamic> map) {
    return BasicUserInfoModel(
      id: map['id'] as String,
      email: map['email'] as String,
      userAvatar:
          map['userAvatar'] != null ? map['userAvatar'] as String : null,
      firstName: map['firstName'] as String,
      lastName: map['lastName'] as String,
    );
  }

  String toJson() => json.encode(toMap());

  factory BasicUserInfoModel.fromJson(String source) =>
      BasicUserInfoModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'BasicUserInfoModel(id: $id, email: $email, userAvatar: $userAvatar, firstName: $firstName, lastName: $lastName)';
  }

  @override
  bool operator ==(covariant BasicUserInfoModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.email == email &&
        other.userAvatar == userAvatar &&
        other.firstName == firstName &&
        other.lastName == lastName;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        email.hashCode ^
        userAvatar.hashCode ^
        firstName.hashCode ^
        lastName.hashCode;
  }
}
