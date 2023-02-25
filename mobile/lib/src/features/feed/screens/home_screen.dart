import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/providers/model_providers.dart';

class HomeScreeen extends ConsumerWidget {
  const HomeScreeen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    return Scaffold(
      body: Column(
        children: [
          SizedBox(
            height: 100,
          ),
          Text(err),
          Center(
            heightFactor: 20,
            child: userInfo != null ? Text(userInfo.email) : Text("nothing"),
          ),
          ElevatedButton(
              onPressed: () {
                ref.watch(myUserProvider.notifier).refetch();
              },
              child: Text("Click me"))
        ],
      ),
    );
  }
}
