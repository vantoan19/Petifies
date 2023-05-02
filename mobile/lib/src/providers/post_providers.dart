import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/uploading_post.dart';

final postInfoProvider =
    Provider<PostModel>((ref) => throw UnimplementedError());

final uploadingPostInfoProvider =
    Provider<UploadingPostModel>((ref) => throw UnimplementedError());
