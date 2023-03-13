// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class UserModel {
  final String id;
  final String email;
  final String? userAvatar;
  final String firstName;
  final String lastName;
  final bool isActivated;
  UserModel({
    required this.id,
    required this.email,
    this.userAvatar = null,
    required this.firstName,
    required this.lastName,
    required this.isActivated,
  });

  UserModel copyWith({
    String? id,
    String? email,
    String? userAvatar,
    String? firstName,
    String? lastName,
    bool? isActivated,
  }) {
    return UserModel(
      id: id ?? this.id,
      email: email ?? this.email,
      userAvatar: userAvatar ?? this.userAvatar,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      isActivated: isActivated ?? this.isActivated,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'email': email,
      'userAvatar': userAvatar,
      'firstName': firstName,
      'lastName': lastName,
      'isActivated': isActivated,
    };
  }

  factory UserModel.fromMap(Map<String, dynamic> map) {
    return UserModel(
      id: map['id'] as String,
      email: map['email'] as String,
      userAvatar:
          map['userAvatar'] != null ? map['userAvatar'] as String : null,
      firstName: map['firstName'] as String,
      lastName: map['lastName'] as String,
      isActivated: map['isActivated'] as bool,
    );
  }

  String toJson() => json.encode(toMap());

  factory UserModel.fromJson(String source) =>
      UserModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'UserModel(id: $id, email: $email, userAvatar: $userAvatar, firstName: $firstName, lastName: $lastName, isActivated: $isActivated)';
  }

  @override
  bool operator ==(covariant UserModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.email == email &&
        other.userAvatar == userAvatar &&
        other.firstName == firstName &&
        other.lastName == lastName &&
        other.isActivated == isActivated;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        email.hashCode ^
        userAvatar.hashCode ^
        firstName.hashCode ^
        lastName.hashCode ^
        isActivated.hashCode;
  }
}
