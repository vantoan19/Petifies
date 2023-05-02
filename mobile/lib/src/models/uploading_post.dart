// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:mobile/src/models/basic_user_info.dart';
import 'package:video_player/video_player.dart';

class UploadingPostModel {
  final String tempId;
  final BasicUserInfoModel owner;
  final String postActivity;
  final DateTime createdAt;
  final String? textContent;
  final List<File>? images;
  final List<VideoPlayerController>? videos;
  UploadingPostModel({
    required this.tempId,
    required this.owner,
    required this.postActivity,
    required this.createdAt,
    this.textContent,
    this.images,
    this.videos,
  });
}
