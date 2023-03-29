// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';

import 'basic_user_info.dart';

class CommentModel {
  final String id;
  final String postID;
  final BasicUserInfoModel owner;
  final DateTime createdAt;
  final String? textContent;
  final NetworkImageModel? image;
  final NetworkVideoModel? video;
  final bool hasReacted;
  final int loveCount;
  final int subcommentCount;
  CommentModel({
    required this.id,
    required this.postID,
    required this.owner,
    required this.createdAt,
    this.textContent,
    this.image,
    this.video,
    required this.hasReacted,
    required this.loveCount,
    required this.subcommentCount,
  });

  CommentModel copyWith({
    String? id,
    String? postID,
    BasicUserInfoModel? owner,
    DateTime? createdAt,
    String? textContent,
    NetworkImageModel? image,
    NetworkVideoModel? video,
    bool? hasReacted,
    int? loveCount,
    int? subcommentCount,
  }) {
    return CommentModel(
      id: id ?? this.id,
      postID: postID ?? this.postID,
      owner: owner ?? this.owner,
      createdAt: createdAt ?? this.createdAt,
      textContent: textContent ?? this.textContent,
      image: image ?? this.image,
      video: video ?? this.video,
      hasReacted: hasReacted ?? this.hasReacted,
      loveCount: loveCount ?? this.loveCount,
      subcommentCount: subcommentCount ?? this.subcommentCount,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'postID': postID,
      'owner': owner.toMap(),
      'createdAt': createdAt.millisecondsSinceEpoch,
      'textContent': textContent,
      'image': image?.toMap(),
      'video': video?.toMap(),
      'hasReacted': hasReacted,
      'loveCount': loveCount,
      'subcommentCount': subcommentCount,
    };
  }

  factory CommentModel.fromMap(Map<String, dynamic> map) {
    return CommentModel(
      id: map['id'] as String,
      postID: map['postID'] as String,
      owner: BasicUserInfoModel.fromMap(map['owner'] as Map<String, dynamic>),
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['createdAt'] as int),
      textContent:
          map['textContent'] != null ? map['textContent'] as String : null,
      image: map['image'] != null
          ? NetworkImageModel.fromMap(map['image'] as Map<String, dynamic>)
          : null,
      video: map['video'] != null
          ? NetworkVideoModel.fromMap(map['video'] as Map<String, dynamic>)
          : null,
      hasReacted: map['hasReacted'] as bool,
      loveCount: map['loveCount'] as int,
      subcommentCount: map['subcommentCount'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory CommentModel.fromJson(String source) =>
      CommentModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'CommentModel(id: $id, postID: $postID, owner: $owner, createdAt: $createdAt, textContent: $textContent, image: $image, video: $video, hasReacted: $hasReacted, loveCount: $loveCount, subcommentCount: $subcommentCount)';
  }

  @override
  bool operator ==(covariant CommentModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.postID == postID &&
        other.owner == owner &&
        other.createdAt == createdAt &&
        other.textContent == textContent &&
        other.image == image &&
        other.video == video &&
        other.hasReacted == hasReacted &&
        other.loveCount == loveCount &&
        other.subcommentCount == subcommentCount;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        postID.hashCode ^
        owner.hashCode ^
        createdAt.hashCode ^
        textContent.hashCode ^
        image.hashCode ^
        video.hashCode ^
        hasReacted.hashCode ^
        loveCount.hashCode ^
        subcommentCount.hashCode;
  }
}
