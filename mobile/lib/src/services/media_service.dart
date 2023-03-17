// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:mobile/src/utils/fileutils.dart';
import 'package:fixnum/fixnum.dart';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/proto/google/protobuf/duration.pb.dart';
import 'package:mobile/src/services/grpc/auth_interceptor.dart';
import 'package:mobile/src/services/grpc/grpc_flutter_client.dart';
import 'package:mobile/src/utils/retryutils.dart';
import 'package:video_player/video_player.dart';

class MediaService {
  AuthGatewayClient? _authClientInstance;
  Ref _ref;

  MediaService({
    required ref,
  }) : this._ref = ref;

  Future<AuthGatewayClient> get _authClient async {
    if (_authClientInstance == null) {
      _authClientInstance = AuthGatewayClient(
          await GrpcFlutterClient.getClient(),
          interceptors: [AuthInterceptor(ref: _ref)]);
    }
    return _authClientInstance!;
  }

  Future<UploadFileResponse> uploadVideo({
    required String uploaderID,
    required File video,
  }) async {
    VideoPlayerController controler = new VideoPlayerController.file(video);
    FileMetadata md = FileMetadata(
      fileName: FileUtils.getFilename(video),
      mediaType: "video",
      uploaderId: uploaderID,
      size: Int64(FileUtils.getFileSize(video)),
      width: controler.value.size.width.toInt(),
      height: controler.value.size.height.toInt(),
      duration: Duration(seconds: Int64(controler.value.duration.inSeconds)),
    );

    Stream<UploadFileRequest> request = Stream.fromIterable(<UploadFileRequest>[
      UploadFileRequest(metadata: md),
      UploadFileRequest(chunkData: video.readAsBytesSync()),
      UploadFileRequest(willBeDiscarded: false),
    ]);

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).userUploadFile(request),
      _ref,
    );
  }

  Future<UploadFileResponse> uploadImage({
    required String uploaderID,
    required File image,
  }) async {
    var decodedImage = await decodeImageFromList(image.readAsBytesSync());
    FileMetadata md = FileMetadata(
      fileName: FileUtils.getFilename(image),
      mediaType: "image",
      uploaderId: uploaderID,
      size: Int64(FileUtils.getFileSize(image)),
      width: decodedImage.width,
      height: decodedImage.height,
    );

    Stream<UploadFileRequest> request = Stream.fromIterable(<UploadFileRequest>[
      UploadFileRequest(metadata: md),
      UploadFileRequest(chunkData: image.readAsBytesSync()),
      UploadFileRequest(willBeDiscarded: false),
    ]);

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).userUploadFile(request),
      _ref,
    );
  }

  Future<RemoveFileByURIResponse> removeFileByURI({
    required String uri,
  }) async {
    RemoveFileByURIRequest request = RemoveFileByURIRequest(uri: uri);

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).removeFileByURI(request),
      _ref,
    );
  }
}
