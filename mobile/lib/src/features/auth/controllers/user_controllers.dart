import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/auth/repositories/user_repository.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/providers/model_providers.dart';

final userCreateControllerProvider =
    AsyncNotifierProvider.autoDispose<UserCreateController, UserModel?>(
        UserCreateController.new);

class UserCreateController extends AutoDisposeAsyncNotifier<UserModel?> {
  @override
  FutureOr<UserModel?> build() {
    return null;
  }

  Future<void> handle({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
  }) async {
    final userRepository = ref.read(userRepositoryProvider);

    state = const AsyncLoading();
    final user = await userRepository.create(
        email: email,
        password: password,
        firstName: firstName,
        lastName: lastName);
    state = user.fold((l) => AsyncValue.error(l, StackTrace.current),
        (r) => AsyncValue.data(r));
  }
}
