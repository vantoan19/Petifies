// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class UserModel {
  final String id;
  final String email;
  final String? userAvatar;
  final String? userWallpaper;
  final String? bio;
  final String firstName;
  final String lastName;
  final bool isActivated;
  final int countPost;
  final int followers;
  final int following;

  UserModel({
    required this.id,
    required this.email,
    this.userAvatar = null,
    this.userWallpaper,
    this.bio,
    required this.firstName,
    required this.lastName,
    required this.isActivated,
    this.countPost = 0,
    this.followers = 0,
    this.following = 0,
  });

  UserModel copyWith({
    String? id,
    String? email,
    String? userAvatar,
    String? userWallpaper,
    String? bio,
    String? firstName,
    String? lastName,
    bool? isActivated,
    int? countPost,
    int? followers,
    int? following,
  }) {
    return UserModel(
      id: id ?? this.id,
      email: email ?? this.email,
      userAvatar: userAvatar ?? this.userAvatar,
      userWallpaper: userWallpaper ?? this.userWallpaper,
      bio: bio ?? this.bio,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      isActivated: isActivated ?? this.isActivated,
      countPost: countPost ?? this.countPost,
      followers: followers ?? this.followers,
      following: following ?? this.following,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'email': email,
      'userAvatar': userAvatar,
      'userWallpaper': userWallpaper,
      'bio': bio,
      'firstName': firstName,
      'lastName': lastName,
      'isActivated': isActivated,
      'countPost': countPost,
      'followers': followers,
      'following': following,
    };
  }

  factory UserModel.fromMap(Map<String, dynamic> map) {
    return UserModel(
      id: map['id'] as String,
      email: map['email'] as String,
      userAvatar:
          map['userAvatar'] != null ? map['userAvatar'] as String : null,
      userWallpaper:
          map['userWallpaper'] != null ? map['userWallpaper'] as String : null,
      bio: map['bio'] != null ? map['bio'] as String : null,
      firstName: map['firstName'] as String,
      lastName: map['lastName'] as String,
      isActivated: map['isActivated'] as bool,
      countPost: map['countPost'] as int,
      followers: map['followers'] as int,
      following: map['following'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory UserModel.fromJson(String source) =>
      UserModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'UserModel(id: $id, email: $email, userAvatar: $userAvatar, userWallpaper: $userWallpaper, bio: $bio, firstName: $firstName, lastName: $lastName, isActivated: $isActivated, countPost: $countPost, followers: $followers, following: $following)';
  }

  @override
  bool operator ==(covariant UserModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.email == email &&
        other.userAvatar == userAvatar &&
        other.userWallpaper == userWallpaper &&
        other.bio == bio &&
        other.firstName == firstName &&
        other.lastName == lastName &&
        other.isActivated == isActivated &&
        other.countPost == countPost &&
        other.followers == followers &&
        other.following == following;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        email.hashCode ^
        userAvatar.hashCode ^
        userWallpaper.hashCode ^
        bio.hashCode ^
        firstName.hashCode ^
        lastName.hashCode ^
        isActivated.hashCode ^
        countPost.hashCode ^
        followers.hashCode ^
        following.hashCode;
  }
}
