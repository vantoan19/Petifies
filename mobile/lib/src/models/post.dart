// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:collection/collection.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/google/protobuf/timestamp.pb.dart';

import 'basic_user_info.dart';

class PostModel {
  final String id;
  final BasicUserInfoModel owner;
  final String postActivity;
  final DateTime createdAt;
  final String? textContent;
  final List<NetworkImageModel>? images;
  final List<NetworkVideoModel>? videos;
  final bool hasReacted;
  final int loveCount;
  final int commentCount;
  PostModel({
    required this.id,
    required this.owner,
    required this.postActivity,
    required this.createdAt,
    this.textContent = null,
    this.images = null,
    this.videos = null,
    required this.hasReacted,
    required this.loveCount,
    required this.commentCount,
  });

  PostModel copyWith({
    String? id,
    BasicUserInfoModel? owner,
    String? postActivity,
    DateTime? createdAt,
    String? textContent,
    List<NetworkImageModel>? images,
    List<NetworkVideoModel>? videos,
    bool? hasReacted,
    int? loveCount,
    int? commentCount,
  }) {
    return PostModel(
      id: id ?? this.id,
      owner: owner ?? this.owner,
      postActivity: postActivity ?? this.postActivity,
      createdAt: createdAt ?? this.createdAt,
      textContent: textContent ?? this.textContent,
      images: images ?? this.images,
      videos: videos ?? this.videos,
      hasReacted: hasReacted ?? this.hasReacted,
      loveCount: loveCount ?? this.loveCount,
      commentCount: commentCount ?? this.commentCount,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'owner': owner.toMap(),
      'postActivity': postActivity,
      'createdAt': createdAt.millisecondsSinceEpoch,
      'textContent': textContent,
      'images': images?.map((x) => x.toMap()).toList(),
      'videos': videos?.map((x) => x.toMap()).toList(),
      'hasReacted': hasReacted,
      'loveCount': loveCount,
      'commentCount': commentCount,
    };
  }

  factory PostModel.fromMap(Map<String, dynamic> map) {
    return PostModel(
      id: map['id'] as String,
      owner: BasicUserInfoModel.fromMap(map['owner'] as Map<String, dynamic>),
      postActivity: map['postActivity'] as String,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['createdAt'] as int),
      textContent:
          map['textContent'] != null ? map['textContent'] as String : null,
      images: map['images'] != null
          ? List<NetworkImageModel>.from(
              (map['images'] as List<int>).map<NetworkImageModel?>(
                (x) => NetworkImageModel.fromMap(x as Map<String, dynamic>),
              ),
            )
          : null,
      videos: map['videos'] != null
          ? List<NetworkVideoModel>.from(
              (map['videos'] as List<int>).map<NetworkVideoModel?>(
                (x) => NetworkVideoModel.fromMap(x as Map<String, dynamic>),
              ),
            )
          : null,
      hasReacted: map['hasReacted'] as bool,
      loveCount: map['loveCount'] as int,
      commentCount: map['commentCount'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory PostModel.fromJson(String source) =>
      PostModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'PostModel(id: $id, owner: $owner, postActivity: $postActivity, createdAt: $createdAt, textContent: $textContent, images: $images, videos: $videos, hasReacted: $hasReacted, loveCount: $loveCount, commentCount: $commentCount)';
  }

  @override
  bool operator ==(covariant PostModel other) {
    if (identical(this, other)) return true;
    final listEquals = const DeepCollectionEquality().equals;

    return other.id == id &&
        other.owner == owner &&
        other.postActivity == postActivity &&
        other.createdAt == createdAt &&
        other.textContent == textContent &&
        listEquals(other.images, images) &&
        listEquals(other.videos, videos) &&
        other.hasReacted == hasReacted &&
        other.loveCount == loveCount &&
        other.commentCount == commentCount;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        owner.hashCode ^
        postActivity.hashCode ^
        createdAt.hashCode ^
        textContent.hashCode ^
        images.hashCode ^
        videos.hashCode ^
        hasReacted.hashCode ^
        loveCount.hashCode ^
        commentCount.hashCode;
  }
}
