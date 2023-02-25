import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/services/user_service.dart';

final userServiceProvider = Provider((ref) => UserService(ref: ref));
