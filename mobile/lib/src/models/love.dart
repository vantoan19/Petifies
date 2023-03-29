// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'basic_user_info.dart';

class LoveModel {
  final String id;
  final BasicUserInfoModel owner;
  final String targetId;
  final bool isPostTarget;
  final DateTime createdAt;
  LoveModel({
    required this.id,
    required this.owner,
    required this.targetId,
    required this.isPostTarget,
    required this.createdAt,
  });

  LoveModel copyWith({
    String? id,
    BasicUserInfoModel? owner,
    String? targetId,
    bool? isPostTarget,
    DateTime? createdAt,
  }) {
    return LoveModel(
      id: id ?? this.id,
      owner: owner ?? this.owner,
      targetId: targetId ?? this.targetId,
      isPostTarget: isPostTarget ?? this.isPostTarget,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'id': id,
      'owner': owner.toMap(),
      'targetId': targetId,
      'isPostTarget': isPostTarget,
      'createdAt': createdAt.millisecondsSinceEpoch,
    };
  }

  factory LoveModel.fromMap(Map<String, dynamic> map) {
    return LoveModel(
      id: map['id'] as String,
      owner: BasicUserInfoModel.fromMap(map['owner'] as Map<String, dynamic>),
      targetId: map['targetId'] as String,
      isPostTarget: map['isPostTarget'] as bool,
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['createdAt'] as int),
    );
  }

  String toJson() => json.encode(toMap());

  factory LoveModel.fromJson(String source) =>
      LoveModel.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'LoveModel(id: $id, owner: $owner, targetId: $targetId, isPostTarget: $isPostTarget, createdAt: $createdAt)';
  }

  @override
  bool operator ==(covariant LoveModel other) {
    if (identical(this, other)) return true;

    return other.id == id &&
        other.owner == owner &&
        other.targetId == targetId &&
        other.isPostTarget == isPostTarget &&
        other.createdAt == createdAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        owner.hashCode ^
        targetId.hashCode ^
        isPostTarget.hashCode ^
        createdAt.hashCode;
  }
}
