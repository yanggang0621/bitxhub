set -e

LOG_FILE_PATH=$1
TARGET=$2
TX_TARGET=$TARGET/tx
mkdir $TX_TARGET

startN=$3
endN=$4
txStartN=$5
txEndN=$6


for (( i=$startN; i<=$endN; i++ ));  do
  echo ===start:$startN, end:$endN, cur:$i

################################ get line
  # a block
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

  ################################ get value
  b_hehght=`sed -n "$verify_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) module.*/\1/g'`
  tx_num=`sed -n "$verify_start p" $LOG_FILE_PATH | head -n 1 | sed 's/.*txNum=\(.*\)/\1/g'`
  # sed -n '22p' ~/work/tmp/bitxhub/scripts/build_solo/logs/bitxhub.log | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'

  # block
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


for (( i=$txStartN; i<=$txEndN; i++ ));  do
  echo ===tx start:$txStartN, tx end:$txEndN, cur:$i

################################ get line
  # a tx
  ## 1. prepare
  tx_prepare_0=`sed -n '/verify prepare start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_prepare_1=`sed -n '/verify prepare end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  ## 2. get
  tx_get_0=`sed -n '/get validator start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_get_1=`sed -n '/get validator end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  ## 3. start
  tx_start_0=`sed -n '/verify start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_start_1=`sed -n '/verify end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  ### 3.1 init
  tx_start_init_0=`sed -n '/init rule start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_start_init_1=`sed -n '/init rule end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  ### 3.2 execute
  tx_start_execute_0=`sed -n '/rule execute start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_start_execute_1=`sed -n '/rule execute end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  ### 3.3 load
  tx_start_load_0=`sed -n '/instances load start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_start_load_1=`sed -n '/instances load end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  ### 3.4 check
  tx_start_check_0=`sed -n '/check status start/=' $LOG_FILE_PATH | head -n $i | tail -n 1`
  tx_start_check_1=`sed -n '/check status end/=' $LOG_FILE_PATH | head -n $i | tail -n 1`

  echo $tx_prepare_0

  ################################ get value
  # tx
  ## 1. prepare
  tx_prepare_0_id=`sed -n "$tx_prepare_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*id=\(.*\) index.*/\1/g'`
  tx_prepare_0_height=`sed -n "$tx_prepare_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) id.*/\1/g'`
  tx_prepare_0_index=`sed -n "$tx_prepare_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_prepare_0_time=`sed -n "$tx_prepare_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_prepare_0_id >> $TX_TARGET/tx_prepare_0.id
  echo $tx_prepare_0_height >> $TX_TARGET/tx_prepare_0.height
  echo $tx_prepare_0_index >> $TX_TARGET/tx_prepare_0.index
  echo $tx_prepare_0_time >> $TX_TARGET/tx_prepare_0.time

  tx_prepare_1_id=`sed -n "$tx_prepare_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*id=\(.*\) index.*/\1/g'`
  tx_prepare_1_height=`sed -n "$tx_prepare_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) id.*/\1/g'`
  tx_prepare_1_index=`sed -n "$tx_prepare_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_prepare_1_time=`sed -n "$tx_prepare_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_prepare_1_id >> $TX_TARGET/tx_prepare_1.id
  echo $tx_prepare_1_height >> $TX_TARGET/tx_prepare_1.height
  echo $tx_prepare_1_index >> $TX_TARGET/tx_prepare_1.index
  echo $tx_prepare_1_time >> $TX_TARGET/tx_prepare_1.time

  ## 2. get
  tx_get_0_height=`sed -n "$tx_get_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_get_0_index=`sed -n "$tx_get_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_get_0_time=`sed -n "$tx_get_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_get_0_height >> $TX_TARGET/tx_get_0.height
  echo $tx_get_0_index >> $TX_TARGET/tx_get_0.index
  echo $tx_get_0_time >> $TX_TARGET/tx_get_0.time

  tx_get_1_height=`sed -n "$tx_get_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_get_1_index=`sed -n "$tx_get_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_get_1_time=`sed -n "$tx_get_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_get_1_height >> $TX_TARGET/tx_get_1.height
  echo $tx_get_1_index >> $TX_TARGET/tx_get_1.index
  echo $tx_get_1_time >> $TX_TARGET/tx_get_1.time

  ## 3. start
  tx_start_0_height=`sed -n "$tx_start_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_0_index=`sed -n "$tx_start_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_0_time=`sed -n "$tx_start_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_0_height >> $TX_TARGET/tx_start_0.height
  echo $tx_start_0_index >> $TX_TARGET/tx_start_0.index
  echo $tx_start_0_time >> $TX_TARGET/tx_start_0.time

  tx_start_1_height=`sed -n "$tx_start_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_1_index=`sed -n "$tx_start_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_1_time=`sed -n "$tx_start_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_1_height >> $TX_TARGET/tx_start_1.height
  echo $tx_start_1_index >> $TX_TARGET/tx_start_1.index
  echo $tx_start_1_time >> $TX_TARGET/tx_start_1.time

  ### 3.1 init
  tx_start_init_0_height=`sed -n "$tx_start_init_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_init_0_index=`sed -n "$tx_start_init_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_init_0_time=`sed -n "$tx_start_init_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_init_0_height >> $TX_TARGET/tx_start_init_0.height
  echo $tx_start_init_0_index >> $TX_TARGET/tx_start_init_0.index
  echo $tx_start_init_0_time >> $TX_TARGET/tx_start_init_0.time

  tx_start_init_1_height=`sed -n "$tx_start_init_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_init_1_index=`sed -n "$tx_start_init_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_init_1_time=`sed -n "$tx_start_init_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_init_1_height >> $TX_TARGET/tx_start_init_1.height
  echo $tx_start_init_1_index >> $TX_TARGET/tx_start_init_1.index
  echo $tx_start_init_1_time >> $TX_TARGET/tx_start_init_1.time

  ### 3.2 execute
  tx_start_execute_0_height=`sed -n "$tx_start_execute_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_execute_0_index=`sed -n "$tx_start_execute_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_execute_0_time=`sed -n "$tx_start_execute_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_execute_0_height >> $TX_TARGET/tx_start_execute_0.height
  echo $tx_start_execute_0_index >> $TX_TARGET/tx_start_execute_0.index
  echo $tx_start_execute_0_time >> $TX_TARGET/tx_start_execute_0.time

  tx_start_execute_1_height=`sed -n "$tx_start_execute_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_execute_1_index=`sed -n "$tx_start_execute_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_execute_1_time=`sed -n "$tx_start_execute_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_execute_1_height >> $TX_TARGET/tx_start_execute_1.height
  echo $tx_start_execute_1_index >> $TX_TARGET/tx_start_execute_1.index
  echo $tx_start_execute_1_time >> $TX_TARGET/tx_start_execute_1.time

  ### 3.3 load
  tx_start_load_0_height=`sed -n "$tx_start_load_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_load_0_index=`sed -n "$tx_start_load_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_load_0_time=`sed -n "$tx_start_load_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_load_0_height >> $TX_TARGET/tx_start_load_0.height
  echo $tx_start_load_0_index >> $TX_TARGET/tx_start_load_0.index
  echo $tx_start_load_0_time >> $TX_TARGET/tx_start_load_0.time

  tx_start_load_1_height=`sed -n "$tx_start_load_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_load_1_index=`sed -n "$tx_start_load_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_load_1_time=`sed -n "$tx_start_load_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_load_1_height >> $TX_TARGET/tx_start_load_1.height
  echo $tx_start_load_1_index >> $TX_TARGET/tx_start_load_1.index
  echo $tx_start_load_1_time >> $TX_TARGET/tx_start_load_1.time

  ### 3.4 check
  tx_start_check_0_height=`sed -n "$tx_start_check_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_check_0_index=`sed -n "$tx_start_check_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_check_0_time=`sed -n "$tx_start_check_0 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_check_0_height >> $TX_TARGET/tx_start_check_0.height
  echo $tx_start_check_0_index >> $TX_TARGET/tx_start_check_0.index
  echo $tx_start_check_0_time >> $TX_TARGET/tx_start_check_0.time

  tx_start_check_1_height=`sed -n "$tx_start_check_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*height=\(.*\) index.*/\1/g'`
  tx_start_check_1_index=`sed -n "$tx_start_check_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*index=\(.*\) module.*/\1/g'`
  tx_start_check_1_time=`sed -n "$tx_start_check_1 p" $LOG_FILE_PATH | head -n 1 | sed 's/.*fields.time=\(.*\) height.*/\1/g'`

  echo $tx_start_check_1_height >> $TX_TARGET/tx_start_check_1.height
  echo $tx_start_check_1_index >> $TX_TARGET/tx_start_check_1.index
  echo $tx_start_check_1_time >> $TX_TARGET/tx_start_check_1.time
done