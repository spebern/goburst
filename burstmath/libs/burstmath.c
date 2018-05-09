// (c) 2018 PoC Consortium ALL RIGHTS RESERVED

#define USE_MULTI_SHABAL
#define _GNU_SOURCE
#define _LARGEFILE64_SOURCE
#include <endian.h>
#include <stddef.h>
#include <string.h>

#include "burstmath.h"
#include "mshabal.h"
#include "mshabal256.h"
#include "shabal.h"

#define SCOOP_SIZE 64
#define NUM_SCOOPS 4096
#define NONCE_SIZE (NUM_SCOOPS * SCOOP_SIZE)

#define HASH_SIZE 32
#define HASH_CAP 4096

#define SET_NONCE(gendata, nonce, offset)                                      \
  xv = (char *)&nonce;                                                         \
  gendata[NONCE_SIZE + offset] = xv[7];                                        \
  gendata[NONCE_SIZE + offset + 1] = xv[6];                                    \
  gendata[NONCE_SIZE + offset + 2] = xv[5];                                    \
  gendata[NONCE_SIZE + offset + 3] = xv[4];                                    \
  gendata[NONCE_SIZE + offset + 4] = xv[3];                                    \
  gendata[NONCE_SIZE + offset + 5] = xv[2];                                    \
  gendata[NONCE_SIZE + offset + 6] = xv[1];                                    \
  gendata[NONCE_SIZE + offset + 7] = xv[0]

uint32_t calculate_scoop(uint64_t height, uint8_t *gensig) {
  shabal_context sc;
  uint8_t new_gensig[32];

  shabal_init(&sc, 256);
  shabal(&sc, gensig, 32);

  uint64_t height_swapped = htobe64(height);
  shabal(&sc, &height_swapped, sizeof(height_swapped));
  shabal_close(&sc, 0, 0, new_gensig);

  return ((new_gensig[30] & 0x0F) << 8) | new_gensig[31];
}

void calculate_deadlines_sse4(CalcDeadlineRequest **reqs, bool poc2){
  char final1[32], final2[32], final3[32], final4[32];
  char gendata1[16 + NONCE_SIZE], gendata2[16 + NONCE_SIZE],
      gendata3[16 + NONCE_SIZE], gendata4[16 + NONCE_SIZE];

  char *xv;

  SET_NONCE(gendata1, reqs[0]->account_id, 0);
  SET_NONCE(gendata2, reqs[1]->account_id, 0);
  SET_NONCE(gendata3, reqs[2]->account_id, 0);
  SET_NONCE(gendata4, reqs[3]->account_id, 0);

  SET_NONCE(gendata1, reqs[0]->nonce, 8);
  SET_NONCE(gendata2, reqs[1]->nonce, 8);
  SET_NONCE(gendata3, reqs[2]->nonce, 8);
  SET_NONCE(gendata4, reqs[3]->nonce, 8);

  mshabal_context x;
  int len;

  for (int i = NONCE_SIZE; i > 0; i -= HASH_SIZE) {
    sse4_mshabal_init(&x, 256);

    len = NONCE_SIZE + 16 - i;
    if (len > HASH_CAP)
      len = HASH_CAP;

    sse4_mshabal(&x, &gendata1[i], &gendata2[i], &gendata3[i], &gendata4[i],
                 len);
    sse4_mshabal_close(&x, 0, 0, 0, 0, 0, &gendata1[i - HASH_SIZE],
                       &gendata2[i - HASH_SIZE], &gendata3[i - HASH_SIZE],
                       &gendata4[i - HASH_SIZE]);
  }

  sse4_mshabal_init(&x, 256);
  sse4_mshabal(&x, gendata1, gendata2, gendata3, gendata4, 16 + NONCE_SIZE);
  sse4_mshabal_close(&x, 0, 0, 0, 0, 0, final1, final2, final3, final4);

  // XOR with final
  for (int i = 0; i < NONCE_SIZE; i++) {
    gendata1[i] ^= (final1[i % 32]);
    gendata2[i] ^= (final2[i % 32]);
    gendata3[i] ^= (final3[i % 32]);
    gendata4[i] ^= (final4[i % 32]);
  }

  uint8_t final11[HASH_SIZE];
  uint8_t final22[HASH_SIZE];
  uint8_t final33[HASH_SIZE];
  uint8_t final44[HASH_SIZE];

  mshabal_context deadline_sc;
  sse4_mshabal_init(&deadline_sc, 256);
  sse4_mshabal(&deadline_sc, reqs[0]->gen_sig, reqs[1]->gen_sig, reqs[2]->gen_sig, reqs[3]->gen_sig, HASH_SIZE);

  uint8_t scoop1[SCOOP_SIZE], scoop2[SCOOP_SIZE], scoop3[SCOOP_SIZE],
      scoop4[SCOOP_SIZE];

  memcpy(scoop1, gendata1 + (reqs[0]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop2, gendata2 + (reqs[1]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop3, gendata3 + (reqs[2]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop4, gendata4 + (reqs[3]->scoop_nr * SCOOP_SIZE), 32);

  if (poc2) {
    memcpy(scoop1 + 32, gendata1 + ((4095 - reqs[0]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop2 + 32, gendata2 + ((4095 - reqs[1]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop3 + 32, gendata3 + ((4095 - reqs[2]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop4 + 32, gendata4 + ((4095 - reqs[3]->scoop_nr) * SCOOP_SIZE) + 32, 32);
  } else {
    memcpy(scoop1 + 32, gendata1 + (reqs[0]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop2 + 32, gendata2 + (reqs[1]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop3 + 32, gendata3 + (reqs[2]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop4 + 32, gendata4 + (reqs[3]->scoop_nr * SCOOP_SIZE) + 32, 32);
  }

  sse4_mshabal(&deadline_sc, scoop1, scoop2, scoop3, scoop4, SCOOP_SIZE);

  sse4_mshabal_close(&deadline_sc, 0, 0, 0, 0, 0, (uint32_t *)final11,
                     (uint32_t *)final22, (uint32_t *)final33,
                     (uint32_t *)final44);

  uint64_t target_result1 = *(uint64_t *)final11;
  uint64_t target_result2 = *(uint64_t *)final22;
  uint64_t target_result3 = *(uint64_t *)final33;
  uint64_t target_result4 = *(uint64_t *)final44;

  *reqs[0]->deadline = target_result1 / reqs[0]->base_target;
  *reqs[1]->deadline = target_result2 / reqs[1]->base_target;
  *reqs[2]->deadline = target_result3 / reqs[2]->base_target;
  *reqs[3]->deadline = target_result4 / reqs[3]->base_target;
}

void calculate_deadlines_avx2(CalcDeadlineRequest **reqs, bool poc2) {
  char final1[32], final2[32], final3[32], final4[32];
  char final5[32], final6[32], final7[32], final8[32];
  char gendata1[16 + NONCE_SIZE], gendata2[16 + NONCE_SIZE],
      gendata3[16 + NONCE_SIZE], gendata4[16 + NONCE_SIZE];
  char gendata5[16 + NONCE_SIZE], gendata6[16 + NONCE_SIZE],
      gendata7[16 + NONCE_SIZE], gendata8[16 + NONCE_SIZE];

  char *xv;

  SET_NONCE(gendata1, reqs[0]->account_id, 0);
  SET_NONCE(gendata2, reqs[1]->account_id, 0);
  SET_NONCE(gendata3, reqs[2]->account_id, 0);
  SET_NONCE(gendata4, reqs[3]->account_id, 0);
  SET_NONCE(gendata5, reqs[4]->account_id, 0);
  SET_NONCE(gendata6, reqs[5]->account_id, 0);
  SET_NONCE(gendata7, reqs[6]->account_id, 0);
  SET_NONCE(gendata8, reqs[7]->account_id, 0);

  SET_NONCE(gendata1, reqs[0]->nonce, 8);
  SET_NONCE(gendata2, reqs[1]->nonce, 8);
  SET_NONCE(gendata3, reqs[2]->nonce, 8);
  SET_NONCE(gendata4, reqs[3]->nonce, 8);
  SET_NONCE(gendata5, reqs[4]->nonce, 8);
  SET_NONCE(gendata6, reqs[5]->nonce, 8);
  SET_NONCE(gendata7, reqs[6]->nonce, 8);
  SET_NONCE(gendata8, reqs[7]->nonce, 8);

  mshabal256_context x;
  int len;

  for (int i = NONCE_SIZE; i;) {
    mshabal256_init(&x);

    len = NONCE_SIZE + 16 - i;
    if (len > HASH_CAP)
      len = HASH_CAP;

    mshabal256(&x, &gendata1[i], &gendata2[i], &gendata3[i], &gendata4[i],
               &gendata5[i], &gendata6[i], &gendata7[i], &gendata8[i], len);

    i -= HASH_SIZE;
    mshabal256_close(&x, (uint32_t *)&gendata1[i], (uint32_t *)&gendata2[i],
                     (uint32_t *)&gendata3[i], (uint32_t *)&gendata4[i],
                     (uint32_t *)&gendata5[i], (uint32_t *)&gendata6[i],
                     (uint32_t *)&gendata7[i], (uint32_t *)&gendata8[i]);
  }

  mshabal256_init(&x);
  mshabal256(&x, gendata1, gendata2, gendata3, gendata4, gendata5, gendata6,
             gendata7, gendata8, 16 + NONCE_SIZE);
  mshabal256_close(&x, (uint32_t *)final1, (uint32_t *)final2,
                   (uint32_t *)final3, (uint32_t *)final4, (uint32_t *)final5,
                   (uint32_t *)final6, (uint32_t *)final7, (uint32_t *)final8);

  // XOR with final
  for (int i = 0; i < NONCE_SIZE; i++) {
    gendata1[i] ^= final1[i % 32];
    gendata2[i] ^= final2[i % 32];
    gendata3[i] ^= final3[i % 32];
    gendata4[i] ^= final4[i % 32];
    gendata5[i] ^= final5[i % 32];
    gendata6[i] ^= final6[i % 32];
    gendata7[i] ^= final7[i % 32];
    gendata8[i] ^= final8[i % 32];
  }

  uint8_t final11[HASH_SIZE];
  uint8_t final22[HASH_SIZE];
  uint8_t final33[HASH_SIZE];
  uint8_t final44[HASH_SIZE];
  uint8_t final55[HASH_SIZE];
  uint8_t final66[HASH_SIZE];
  uint8_t final77[HASH_SIZE];
  uint8_t final88[HASH_SIZE];

  mshabal256_context deadline_sc;
  mshabal256_init(&deadline_sc);
  mshabal256(&deadline_sc, reqs[0]->gen_sig, reqs[1]->gen_sig, reqs[2]->gen_sig, reqs[3]->gen_sig, reqs[4]->gen_sig,
             reqs[5]->gen_sig, reqs[6]->gen_sig, reqs[7]->gen_sig, HASH_SIZE);

  uint8_t scoop1[SCOOP_SIZE], scoop2[SCOOP_SIZE], scoop3[SCOOP_SIZE],
      scoop4[SCOOP_SIZE], scoop5[SCOOP_SIZE], scoop6[SCOOP_SIZE],
      scoop7[SCOOP_SIZE], scoop8[SCOOP_SIZE];

  memcpy(scoop1, gendata1 + (reqs[0]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop2, gendata2 + (reqs[1]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop3, gendata3 + (reqs[2]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop4, gendata4 + (reqs[3]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop5, gendata5 + (reqs[4]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop6, gendata6 + (reqs[5]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop7, gendata7 + (reqs[6]->scoop_nr * SCOOP_SIZE), 32);
  memcpy(scoop8, gendata8 + (reqs[7]->scoop_nr * SCOOP_SIZE), 32);

  if (poc2) {
    memcpy(scoop1 + 32, gendata1 + ((4095 - reqs[0]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop2 + 32, gendata2 + ((4095 - reqs[1]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop3 + 32, gendata3 + ((4095 - reqs[2]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop4 + 32, gendata4 + ((4095 - reqs[3]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop5 + 32, gendata5 + ((4095 - reqs[4]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop6 + 32, gendata6 + ((4095 - reqs[5]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop7 + 32, gendata7 + ((4095 - reqs[6]->scoop_nr) * SCOOP_SIZE) + 32, 32);
    memcpy(scoop8 + 32, gendata8 + ((4095 - reqs[7]->scoop_nr) * SCOOP_SIZE) + 32, 32);
  } else {
    memcpy(scoop1 + 32, gendata1 + (reqs[0]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop2 + 32, gendata2 + (reqs[1]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop3 + 32, gendata3 + (reqs[2]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop4 + 32, gendata4 + (reqs[3]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop5 + 32, gendata5 + (reqs[4]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop6 + 32, gendata6 + (reqs[5]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop7 + 32, gendata7 + (reqs[6]->scoop_nr * SCOOP_SIZE) + 32, 32);
    memcpy(scoop8 + 32, gendata8 + (reqs[7]->scoop_nr * SCOOP_SIZE) + 32, 32);
  }

  mshabal256(&deadline_sc, scoop1, scoop2, scoop3, scoop4, scoop5, scoop6,
             scoop7, scoop8, SCOOP_SIZE);

  mshabal256_close(&deadline_sc, (uint32_t *)final11, (uint32_t *)final22,
                   (uint32_t *)final33, (uint32_t *)final44,
                   (uint32_t *)final55, (uint32_t *)final66,
                   (uint32_t *)final77, (uint32_t *)final88);

  uint64_t target_result1 = *(uint64_t *)final11;
  uint64_t target_result2 = *(uint64_t *)final22;
  uint64_t target_result3 = *(uint64_t *)final33;
  uint64_t target_result4 = *(uint64_t *)final44;
  uint64_t target_result5 = *(uint64_t *)final55;
  uint64_t target_result6 = *(uint64_t *)final66;
  uint64_t target_result7 = *(uint64_t *)final77;
  uint64_t target_result8 = *(uint64_t *)final88;

  *reqs[0]->deadline = target_result1 / reqs[0]->base_target;
  *reqs[1]->deadline = target_result2 / reqs[1]->base_target;
  *reqs[2]->deadline = target_result3 / reqs[2]->base_target;
  *reqs[3]->deadline = target_result4 / reqs[3]->base_target;
  *reqs[4]->deadline = target_result5 / reqs[4]->base_target;
  *reqs[5]->deadline = target_result6 / reqs[5]->base_target;
  *reqs[6]->deadline = target_result7 / reqs[6]->base_target;
  *reqs[7]->deadline = target_result8 / reqs[7]->base_target;
}
