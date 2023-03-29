// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:mobile/src/models/basic_user_info.dart';
import 'package:video_player/video_player.dart';

class UploadingCommentModel {
  final String tempId;
  final BasicUserInfoModel owner;
  final DateTime createdAt;
  final String? textContent;
  final File? image;
  final VideoPlayerController? video;
  UploadingCommentModel({
    required this.tempId,
    required this.owner,
    required this.createdAt,
    this.textContent,
    this.image,
    this.video,
  });
}
