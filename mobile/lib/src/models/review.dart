// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';

class ReviewModel {
  final String id;
  final String petifiesId;
  final BasicUserInfoModel author;
  final String review;
  final NetworkImageModel? image;
  final DateTime createdAt;
  ReviewModel({
    required this.id,
    required this.petifiesId,
    required this.author,
    required this.review,
    this.image,
    required this.createdAt,
  });

  ReviewModel copyWith({
    String? id,
    String? petifiesId,
    BasicUserInfoModel? author,
    String? review,
    NetworkImageModel? image,
    DateTime? createdAt,
  }) {
    return ReviewModel(
      id: id ?? this.id,
      petifiesId: petifiesId ?? this.petifiesId,
      author: author ?? this.author,
      review: review ?? this.review,
      image: image ?? this.image,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'petifiesId': petifiesId,
      'author': author.toMap(),
      'review': review,
      'image': image?.toMap(),
      'createdAt': createdAt.millisecondsSinceEpoch,
    };
  }

  factory ReviewModel.fromMap(Map<String, dynamic> map) {
    return ReviewModel(
      id: map['id'] as String,
      petifiesId: map['petifiesId'] as String,
      author: BasicUserInfoModel.fromMap(map['author'] as Map<String, dynamic>),
      review: map['review'] as String,
      image: map['image'] != null
          ? NetworkImageModel.fromMap(map['image'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['createdAt'] as int),
    );
  }

  String toJson() => json.encode(toMap());

  factory ReviewModel.fromJson(String source) =>
      ReviewModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'Review(id: $id, petifiesId: $petifiesId, author: $author, review: $review, image: $image, createdAt: $createdAt)';
  }

  @override
  bool operator ==(covariant ReviewModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.petifiesId == petifiesId &&
        other.author == author &&
        other.review == review &&
        other.image == image &&
        other.createdAt == createdAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        petifiesId.hashCode ^
        author.hashCode ^
        review.hashCode ^
        image.hashCode ^
        createdAt.hashCode;
  }
}
