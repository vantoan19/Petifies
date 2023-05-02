import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/controllers/petifies_controllers.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/widgets/appbars/my_petifies_appbar.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';
import 'package:mobile/src/widgets/petifies/petifies.dart';

class MyPetifiesScreen extends ConsumerWidget {
  const MyPetifiesScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);
    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => null,
    );

    if (err != null) {
      return Center(child: Text(err));
    }

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    if (userInfo == null) {
      return Center(child: CircularProgressIndicator());
    }

    final petifies =
        ref.watch(listPetifiesByUserIdProvider(userId: userInfo.id));

    return Scaffold(
      appBar: OnlyGoBackAppbar(),
      body: petifies.when(
        data: (data) {
          return Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: Constants.petifiesExpoloreHorizontalPadding,
                  vertical: 8,
                ),
                child: Text(
                  "My Petifies",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
              Expanded(
                child: ListView.builder(
                  itemBuilder: (context, index) => ProviderScope(
                    child: ProviderScope(
                      key: ObjectKey(data[index].id),
                      overrides: [
                        petifiesInfoProvider.overrideWithValue(data[index])
                      ],
                      child: const Petifies(),
                    ),
                  ),
                  itemCount: data.length,
                ),
              ),
            ],
          );
        },
        error: (error, stackTrace) => Center(
          child: Text(error.toString()),
        ),
        loading: () => Center(
          child: CircularProgressIndicator(),
        ),
      ),
    );
  }
}
