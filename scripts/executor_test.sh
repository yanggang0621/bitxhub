set -e

LOG_FILE_PATH=$1
TARGET=$2

startN=$3
endN=$4

for (( i=$startN; i<=$endN; i++ ));  do
  echo ===start:$startN, end:$endN, cur:$i

  proof_start=`sed -n '/check proof start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  proof_end=`sed -n '/check proof end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  # sed -n '/verify proofs 【start】/=' ~/work/tmp/bitxhub/scripts/build_solo/logs/bitxhub.log | head -n 1
  verify_start=`sed -n '/verify proofs 【start】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  verify_end=`sed -n '/verify proofs 【end】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  apply_start=`sed -n '/apply transaction 【start】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  apply_end=`sed -n '/apply transaction 【end】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  build_start=`sed -n '/build tx merkle tree 【start】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  build_end=`sed -n '/build tx merkle tree 【end】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  calc_start=`sed -n '/calc receipt merkle root 【start】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  calc_end=`sed -n '/calc receipt merkle root 【end】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  calc_timeout_start=`sed -n '/calc timeout l2 root 【start】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  calc_timeout_end=`sed -n '/calc timeout l2 root 【end】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  timeout_root_start=`sed -n '/calc merkle root 【start】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  timeout_root_end=`sed -n '/calc merkle root 【end】/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  echo $verify_start
  b_hehght=`sed -n "$verify_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) module.*/\1/g'`
  tx_num=`sed -n "$verify_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*txNum=\(.*\)/\1/g'`
  # sed -n '22p' ~/work/tmp/bitxhub/scripts/build_solo/logs/bitxhub.log | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'

  proof_start_time=`sed -n "$proof_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  proof_end_time=`sed -n "$proof_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  verify_start_time=`sed -n "$verify_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  verify_end_time=`sed -n "$verify_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  apply_start_time=`sed -n "$apply_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  apply_end_time=`sed -n "$apply_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  build_start_time=`sed -n "$build_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  build_end_time=`sed -n "$build_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  calc_start_time=`sed -n "$calc_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  calc_end_time=`sed -n "$calc_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  calc_timeout_start_time=`sed -n "$calc_timeout_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  calc_timeout_end_time=`sed -n "$calc_timeout_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  timeout_root_start_time=`sed -n "$timeout_root_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`
  timeout_root_end_time=`sed -n "$timeout_root_end p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $b_hehght >> $TARGET/block.height
  echo $tx_num >> $TARGET/tx.num

  echo $proof_start_time >> $TARGET/proof_start.time
  echo $proof_end_time >> $TARGET/proof_end.time

  echo $verify_start_time >> $TARGET/verify_start.time
  echo $verify_end_time >> $TARGET/verify_end.time

  echo $apply_start_time >> $TARGET/apply_star.time
  echo $apply_end_time >> $TARGET/apply_end.time

  echo $build_start_time >> $TARGET/build_start.time
  echo $build_end_time >> $TARGET/build_end.time

  echo $calc_start_time >> $TARGET/calc_start.time
  echo $calc_end_time >> $TARGET/calc_end.time

  echo $calc_timeout_start_time >> $TARGET/calc_timeout_start.time
  echo $calc_timeout_end_time >> $TARGET/calc_timeout_end.time

  echo $timeout_root_start_time >> $TARGET/timeout_root_start.time
  echo $timeout_root_end_time >> $TARGET/timeout_root_end.time

done