/**
 * Copyright 2017 ~ 2025 the original author or authors[983708408@qq.com].
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this export except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"github.com/Shopify/sarama"
	"github.com/wl4g/kafka_offset_tool/pkg/common"
	"log"
	"sync"
)

type ProducedOffset struct {
	NewestOffset int64 // LogSize
	OldestOffset int64
}

/**
 * Produced topic partition offsets.
 * @author Wang.sir <wanglsir@gmail.com,983708408@qq.com>
 * @date 19-07-20
 */
func getProducedTopicPartitionOffsets() map[string]map[int32]ProducedOffset {
	log.Printf("Fetching metadata of the topic partitions infor...")

	// Describe topic partition offset.
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	producedTopicOffsets := make(map[string]map[int32]ProducedOffset)
	for _, topic := range listTopicAll() {
		//log.Printf("Fetching partition info for topics: %s ...", topic)

		go func(topic string) {
			wg.Add(1)
			defer wg.Done()
			partitions, err := opt.client.Partitions(topic)
			if err != nil {
				common.ErrorExit(err, "Cannot get partitions of topic: %s, %s", topic)
			}
			mu.Lock()
			producedTopicOffsets[topic] = make(map[int32]ProducedOffset, len(partitions))
			mu.Unlock()

			for _, partition := range partitions {
				//fmt.Printf("topic:%s, part:%d \n", topic, partition)
				mu.Lock()
				_topicOffset := ProducedOffset{}
				mu.Unlock()

				// Largest offset(logSize).
				newestOffset, err := opt.client.GetOffset(topic, partition, sarama.OffsetNewest)
				if err != nil {
					common.ErrorExit(err, "Cannot get current offset of topic: %s, partition: %d, %s", topic, partition)
				} else {
					mu.Lock()
					_topicOffset.NewestOffset = newestOffset
					mu.Unlock()
				}

				// Oldest offset.
				oldestOffset, err := opt.client.GetOffset(topic, partition, sarama.OffsetOldest)
				if err != nil {
					common.ErrorExit(err, "Cannot get current oldest offset of topic: %s, partition: %d, %s", topic, partition)
				} else {
					mu.Lock()
					_topicOffset.OldestOffset = oldestOffset
					mu.Unlock()
				}

				mu.Lock()
				producedTopicOffsets[topic][partition] = _topicOffset
				mu.Unlock()
			}
		}(topic)
	}
	wg.Wait()
	return producedTopicOffsets
}
