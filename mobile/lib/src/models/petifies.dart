// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:collection/collection.dart';
import 'package:mobile/src/models/address.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';

enum PetifiesType {
  DOG_WALKING,
  CAT_PLAYING,
  DOG_SITTING,
  CAT_SITTING,
  DOG_ADOPTION,
  CAT_ADOPTION,
  UNKNOWN,
}

enum PetifiesStatus {
  UNAVAILABLE,
  AVAILABLE,
  DELETED,
  UNKNOWN,
}

class PetifiesModel {
  final String id;
  final BasicUserInfoModel owner;
  final PetifiesType type;
  final String title;
  final String description;
  final String petName;
  final List<NetworkImageModel> image;
  final PetifiesStatus status;
  final Address address;
  final DateTime createdAt;
  PetifiesModel({
    required this.id,
    required this.owner,
    required this.type,
    required this.title,
    required this.description,
    required this.petName,
    required this.image,
    required this.status,
    required this.address,
    required this.createdAt,
  });

  PetifiesModel copyWith({
    String? id,
    BasicUserInfoModel? owner,
    PetifiesType? type,
    String? title,
    String? description,
    String? petName,
    List<NetworkImageModel>? image,
    PetifiesStatus? status,
    Address? address,
    DateTime? createdAt,
  }) {
    return PetifiesModel(
      id: id ?? this.id,
      owner: owner ?? this.owner,
      type: type ?? this.type,
      title: title ?? this.title,
      description: description ?? this.description,
      petName: petName ?? this.petName,
      image: image ?? this.image,
      status: status ?? this.status,
      address: address ?? this.address,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  bool operator ==(covariant PetifiesModel other) {
    if (identical(this, other)) return true;
    final listEquals = const DeepCollectionEquality().equals;

    return other.id == id &&
        other.owner == owner &&
        other.type == type &&
        other.title == title &&
        other.description == description &&
        other.petName == petName &&
        listEquals(other.image, image) &&
        other.status == status &&
        other.address == address &&
        other.createdAt == createdAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        owner.hashCode ^
        type.hashCode ^
        title.hashCode ^
        description.hashCode ^
        petName.hashCode ^
        image.hashCode ^
        status.hashCode ^
        address.hashCode ^
        createdAt.hashCode;
  }

  @override
  String toString() {
    return 'PetifiesModel(id: $id, owner: $owner, type: $type, title: $title, description: $description, petName: $petName, image: $image, status: $status, address: $address, createdAt: $createdAt)';
  }
}
