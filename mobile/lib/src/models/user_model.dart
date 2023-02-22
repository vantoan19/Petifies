// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class UserModel {
  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final bool isAuthenticated;
  final bool isActivated;
  UserModel({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.isAuthenticated,
    required this.isActivated,
  });

  UserModel copyWith({
    String? id,
    String? email,
    String? firstName,
    String? lastName,
    bool? isAuthenticated,
    bool? isActivated,
  }) {
    return UserModel(
      id: id ?? this.id,
      email: email ?? this.email,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      isActivated: isActivated ?? this.isActivated,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'email': email,
      'firstName': firstName,
      'lastName': lastName,
      'isAuthenticated': isAuthenticated,
      'isActivated': isActivated,
    };
  }

  factory UserModel.fromMap(Map<String, dynamic> map) {
    return UserModel(
      id: map['id'] as String,
      email: map['email'] as String,
      firstName: map['firstName'] as String,
      lastName: map['lastName'] as String,
      isAuthenticated: map['isAuthenticated'] as bool,
      isActivated: map['isActivated'] as bool,
    );
  }

  String toJson() => json.encode(toMap());

  factory UserModel.fromJson(String source) =>
      UserModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'UserModel(id: $id, email: $email, firstName: $firstName, lastName: $lastName, isAuthenticated: $isAuthenticated, isActivated: $isActivated)';
  }

  @override
  bool operator ==(covariant UserModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.email == email &&
        other.firstName == firstName &&
        other.lastName == lastName &&
        other.isAuthenticated == isAuthenticated &&
        other.isActivated == isActivated;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        email.hashCode ^
        firstName.hashCode ^
        lastName.hashCode ^
        isAuthenticated.hashCode ^
        isActivated.hashCode;
  }
}
