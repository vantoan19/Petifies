import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/auth/repositories/user_repository.dart';
import 'package:mobile/src/models/user_model.dart';

final myUserProvider =
    AsyncNotifierProvider.autoDispose<MyUser, UserModel?>(MyUser.new);

class MyUser extends AutoDisposeAsyncNotifier<UserModel?> {
  @override
  FutureOr<UserModel?> build() async {
    final userRepository = ref.read(userRepositoryProvider);

    final user = await userRepository.getMyInfo();
    return user.fold((l) => null, (r) => r);
  }

  void SetUser(UserModel? user) {
    state = AsyncValue.data(user);
  }

  Future<void> refetch() async {
    final userRepository = ref.read(userRepositoryProvider);

    state = AsyncLoading();
    final result = await userRepository.getMyInfo();
    state = result.fold((l) => AsyncValue.error(l, StackTrace.current),
        (r) => AsyncValue.data(r));
  }
}
