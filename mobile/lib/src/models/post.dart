// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:collection/collection.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/google/protobuf/timestamp.pb.dart';

class PostModel {
  final UserModel owner;
  final String postActivity;
  final Timestamp postTime;
  final String? textContent;
  final List<NetworkImageModel>? images;
  final List<NetworkVideoModel>? videos;
  final int loveCount;
  final int commentCount;
  PostModel({
    required this.owner,
    required this.postActivity,
    required this.postTime,
    this.textContent = null,
    this.images = null,
    this.videos = null,
    required this.loveCount,
    required this.commentCount,
  });

  PostModel copyWith({
    UserModel? owner,
    String? postActivity,
    Timestamp? postTime,
    String? textContent,
    List<NetworkImageModel>? images,
    List<NetworkVideoModel>? videos,
    int? loveCount,
    int? commentCount,
  }) {
    return PostModel(
      owner: owner ?? this.owner,
      postActivity: postActivity ?? this.postActivity,
      postTime: postTime ?? this.postTime,
      textContent: textContent ?? this.textContent,
      images: images ?? this.images,
      videos: videos ?? this.videos,
      loveCount: loveCount ?? this.loveCount,
      commentCount: commentCount ?? this.commentCount,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'owner': owner.toMap(),
      'postActivity': postActivity,
      'postTime': postTime,
      'textContent': textContent,
      'images': images?.map((x) => x.toMap()).toList(),
      'videos': videos?.map((x) => x.toMap()).toList(),
      'loveCount': loveCount,
      'commentCount': commentCount,
    };
  }

  factory PostModel.fromMap(Map<String, dynamic> map) {
    return PostModel(
      owner: UserModel.fromMap(map['owner'] as Map<String, dynamic>),
      postActivity: map['postActivity'] as String,
      postTime: map['postTime'],
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
      loveCount: map['loveCount'] as int,
      commentCount: map['commentCount'] as int,
    );
  }

  String toJson() => json.encode(toMap());

  factory PostModel.fromJson(String source) =>
      PostModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'PostModel(owner: $owner, postActivity: $postActivity, postTime: $postTime, textContent: $textContent, images: $images, videos: $videos, loveCount: $loveCount, commentCount: $commentCount)';
  }

  @override
  bool operator ==(covariant PostModel other) {
    if (identical(this, other)) return true;
    final listEquals = const DeepCollectionEquality().equals;

    return other.owner == owner &&
        other.postActivity == postActivity &&
        other.postTime == postTime &&
        other.textContent == textContent &&
        listEquals(other.images, images) &&
        listEquals(other.videos, videos) &&
        other.loveCount == loveCount &&
        other.commentCount == commentCount;
  }

  @override
  int get hashCode {
    return owner.hashCode ^
        postActivity.hashCode ^
        postTime.hashCode ^
        textContent.hashCode ^
        images.hashCode ^
        videos.hashCode ^
        loveCount.hashCode ^
        commentCount.hashCode;
  }
}
