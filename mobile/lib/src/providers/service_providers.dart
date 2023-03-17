import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/services/media_service.dart';
import 'package:mobile/src/services/post_service.dart';
import 'package:mobile/src/services/user_service.dart';

final userServiceProvider = Provider((ref) => UserService(ref: ref));
final postServiceProvider = Provider((ref) => PostService(ref: ref));
final mediaServiceProvider = Provider((ref) => MediaService(ref: ref));
