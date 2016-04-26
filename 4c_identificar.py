# These are all the modules we'll be using later. Make sure you can import them
# before proceeding further.
from __future__ import print_function
import numpy as np
import tensorflow as tf
import sys
from scipy import ndimage
from six.moves import cPickle as pickle
from six.moves import range

#Reformat into a TensorFlow-friendly shape:
#	- convolutions need the image data formatted as a cube (width by height by #channels)
#	- labels as float 1-hot encodings.
image_size = 28
num_labels = 10
num_channels = 1 # grayscale
pixel_depth = 255.0

import numpy as np


def load_letter():
	dataset = np.ndarray(shape=(1, image_size, image_size), dtype=np.float32)
	image_index = 0
	image_file = "char.png"
	try:
		image_data = (ndimage.imread(image_file).astype(float) - pixel_depth / 2) / pixel_depth
		if image_data.shape != (image_size, image_size):
			raise Exception('Unexpected image shape: %s' % str(image_data.shape))
	except IOError as e:
		print('Could not read:', image_file, ':', e, '- it\'s ok, skipping.')
	return image_data

	
letras = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J']

with open('trained.bin', 'rb') as f:
	entrenado = pickle.load(f)
	# Cargar datos entrenados
	layer1_weights = entrenado[0]
	layer1_biases = entrenado[1]
	layer2_weights = entrenado[2]
	layer2_biases = entrenado[3]
	layer3_weights = entrenado[4]
	layer3_biases = entrenado[5]
	layer4_weights = entrenado[6]
	layer4_biases = entrenado[7]
	problem_dataset = load_letter()

  

def reformat(dataset, labels):
  dataset = dataset.reshape(
    (-1, image_size, image_size, num_channels)).astype(np.float32)
  labels = (np.arange(num_labels) == labels[:,None]).astype(np.float32)
  return dataset, labels

problem_dataset = problem_dataset.reshape((-1, image_size, image_size, num_channels)).astype(np.float32)

def accuracy(predictions, labels):
  return (100.0 * np.sum(np.argmax(predictions, 1) == np.argmax(labels, 1))
          / predictions.shape[0])
		  
# Let's build a small network with two convolutional layers, followed by one fully connected layer.
# Convolutional networks are more expensive computationally, so we'll limit its depth and number of
# fully connected nodes.
batch_size = 16
patch_size = 5
depth = 16
num_hidden = 64

graph = tf.Graph()

with graph.as_default():

  # Input data.
  tf_problem_dataset = tf.constant(problem_dataset)
  
  # Model.
  def model(data):
    conv = tf.nn.conv2d(data, layer1_weights, [1, 2, 2, 1], padding='SAME')
    hidden = tf.nn.relu(conv + layer1_biases)
    conv = tf.nn.conv2d(hidden, layer2_weights, [1, 2, 2, 1], padding='SAME')
    hidden = tf.nn.relu(conv + layer2_biases)
    shape = hidden.get_shape().as_list()
    reshape = tf.reshape(hidden, [shape[0], shape[1] * shape[2] * shape[3]])
    hidden = tf.nn.relu(tf.matmul(reshape, layer3_weights) + layer3_biases)
    return tf.matmul(hidden, layer4_weights) + layer4_biases
  
  # Training computation.
  #logits = model(tf_train_dataset)
  #loss = tf.reduce_mean(tf.nn.softmax_cross_entropy_with_logits(logits, tf_train_labels))
    
  # Optimizer.
  #optimizer = tf.train.GradientDescentOptimizer(0.05).minimize(loss)
  
  # Prediction:
  problem_prediction = tf.nn.softmax(model(tf_problem_dataset))
  
  
num_steps = 1001

with tf.Session(graph=graph) as session:
  tf.initialize_all_variables().run()
  resultado = np.array(problem_prediction.eval())
  #print('Con probabilidad %f%% la letra es una: ' % resultado.max())
  print(letras[resultado.argmax()])
	
#tf.nn.max_pool([batch, height, width, channels], [4,4,4,4], [3,3,3,3], 'VALID', data_format='NHWC', name=None)